package server

import (
	"fmt"
	"os"
	"sync"
	"time"

	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/mq"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

// Server interface
type Server interface {
	StartBuilding(target target.Target, pattern pattern.Pattern, params interface{}) error
	StopBuilding(target target.Target, pattern pattern.Pattern) error
	GetResult(target target.Target, pattern pattern.Pattern, from, to uint) (interface{}, error)
	Close() error
}

type serverImpl struct {
	impl           builder.Implement
	mq             mq.MQ
	workers        map[string]worker
	consumedCounts map[string]uint64
	createWorker   func(string, interface{}, builder.Builder, mq.Codec) worker
	logger         *logger.Logger
	locker         sync.Locker
	isClosed       bool
	maxRetainCount uint64
	monitorPeriod  time.Duration
}

// worker unique id
func workerUID(target target.Target, pattern pattern.Pattern) string {
	return string(target) + "_" + string(pattern)
}

// StartBuilding the building
func (s *serverImpl) StartBuilding(target target.Target,
	pattern pattern.Pattern,
	params interface{},
) error {
	codec := mq.GetCodec(target, pattern)
	if codec == nil {
		return fmt.Errorf("no codec of target:%s,pattern:%s", target, pattern)
	}

	logger.Debugf("start build %s:%s", target, pattern)
	dataBuilder := builder.GetBuilder(s.impl, target, pattern)
	if dataBuilder == nil {
		return fmt.Errorf(
			"no build implemented of impl:%s, target:%s, pattern:%s",
			s.impl,
			target,
			pattern)
	}

	uid := workerUID(target, pattern)
	s.locker.Lock()
	if s.workers[uid] != nil {
		s.locker.Unlock()
		return fmt.Errorf("worker of target:%s, pattern:%s exits",
			target,
			pattern)
	}
	worker := s.createWorker(uid, params, dataBuilder, codec)
	s.workers[uid] = worker
	s.consumedCounts[uid] = 0
	s.locker.Unlock()
	// make sure queue is clean
	s.mq.DeleteTopic(uid)
	worker.start()
	return nil
}

// StopBuilding the building
func (s *serverImpl) StopBuilding(target target.Target,
	pattern pattern.Pattern,
) error {
	logger.Debugf("start build %s:%s", target, pattern)
	uid := workerUID(target, pattern)
	s.locker.Lock()
	worker := s.workers[uid]
	if worker == nil {
		s.locker.Unlock()
		return fmt.Errorf("no worker exists of target:%s, pattern:%s",
			target, pattern)
	}
	delete(s.workers, uid)
	delete(s.consumedCounts, uid)
	s.locker.Unlock()
	worker.stop()
	return nil
}

// GetResult returns the building result
func (s *serverImpl) GetResult(target target.Target,
	pattern pattern.Pattern,
	from, to uint,
) (interface{}, error) {
	codec := mq.GetCodec(target, pattern)
	if codec == nil {
		return nil, fmt.Errorf("no codec of target:%s,pattern:%s", target, pattern)
	}

	uid := workerUID(target, pattern)
	s.locker.Lock()
	if s.workers[uid] == nil {
		s.locker.Unlock()
		return nil, fmt.Errorf("no worker exists of target:%s, pattern:%s",
			target, pattern)
	}
	s.locker.Unlock()

	msgs, err := s.mq.Get(uid, from, to, codec)
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, len(msgs))
	for i, m := range msgs {
		ret[i] = m.Body
	}

	s.locker.Lock()
	w := s.workers[uid]
	if w == nil {
		s.logger.Errorf("%s is stopped", uid)
	} else {
		c := uint64(len(ret)) + s.consumedCounts[uid]
		s.consumedCounts[uid] = c
		stat := w.status()
		if stat.status() == Pause && stat.buildCount() < c+s.maxRetainCount {
			w.proceed()
		}
	}
	s.locker.Unlock()

	return ret, nil
}

// Close workers
func (s *serverImpl) Close() error {
	for uid, w := range s.workers {
		_ = uid
		s.logger.Infof("%s status %d", uid, w.status().status())
		w.stop()
	}
	s.locker.Lock()
	s.workers = nil
	s.consumedCounts = nil
	s.locker.Unlock()
	return nil
}

func (s *serverImpl) monitor() {
	for !s.isClosed {
		select {
		case <-time.After(s.monitorPeriod):
		}

		s.locker.Lock()
		for _, w := range s.workers {
			name := w.name()
			c := s.consumedCounts[name]
			stat := w.status()
			curStat, bc := stat.status(), stat.buildCount()
			switch curStat {
			case Start:
				if bc > c+s.maxRetainCount {
					s.logger.Infof("pause worker:%s, consumed:%d, buildCount:%d",
						name, c, bc)
					w.pause()
				}
			case Pause:
				if bc < c+s.maxRetainCount {
					w.proceed()
				}
			}
		}
		s.locker.Unlock()
	}
}

// CreateServer create build server
func CreateServer(impl builder.Implement, q mq.MQ) (Server, error) {
	if !builder.HasImplement(impl) {
		return nil, fmt.Errorf("no implementation of %s", impl)
	}

	if q == nil {
		return nil, fmt.Errorf("nil mq")
	}

	srv := &serverImpl{
		impl:           impl,
		mq:             q,
		workers:        make(map[string]worker),
		consumedCounts: make(map[string]uint64),
		logger:         logger.New(os.Stderr, "[server] ", logger.Ldefault),
		locker:         new(sync.Mutex),
		isClosed:       false,
		maxRetainCount: uint64(1000000),
		monitorPeriod:  time.Second * 3,
	}
	srv.logger.Level = logger.Ldebug

	srv.createWorker = func(uid string, params interface{},
		dataBuilder builder.Builder, codec mq.Codec) worker {
		return &workerImpl{
			uid:         uid,
			mq:          q,
			codec:       codec,
			params:      params,
			dataBuilder: dataBuilder,
			logger:      srv.logger,
		}
	}

	go srv.monitor()
	return srv, nil
}
