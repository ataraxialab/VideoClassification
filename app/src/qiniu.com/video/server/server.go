package server

import (
	"fmt"

	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/mq"
)

// Server interface
type Server interface {
	StartBuilding(target builder.Target, pattern builder.Pattern, params interface{}) error
	StopBuilding(target builder.Target, pattern builder.Pattern) error
	GetResult(target builder.Target, pattern builder.Pattern, from, to uint) (interface{}, error)
	Close() error
}

type serverImpl struct {
	impl         builder.Implement
	mq           mq.MQ
	workers      map[string]worker
	createWorker func(string, interface{}, builder.Builder, mq.Codec) worker
	logger       *logger.Logger
}

// worker unique id
func workerUID(target builder.Target, pattern builder.Pattern) string {
	return string(target) + "_" + string(pattern)
}

// StartBuilding the building
func (s *serverImpl) StartBuilding(target builder.Target,
	pattern builder.Pattern,
	params interface{},
) error {
	codec := mq.GetCodec(target, pattern)
	if codec == nil {
		return fmt.Errorf("no codec of target:%s,pattern:%s", target, pattern)
	}

	uid := workerUID(target, pattern)
	if s.workers[uid] != nil {
		return fmt.Errorf("worker of target:%s, pattern:%s exits",
			target,
			pattern)
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
	worker := s.createWorker(uid, params, dataBuilder, codec)
	s.workers[uid] = worker
	worker.start()
	return nil
}

// StopBuilding the building
func (s *serverImpl) StopBuilding(target builder.Target,
	pattern builder.Pattern,
) error {
	logger.Debugf("start build %s:%s", target, pattern)
	uid := workerUID(target, pattern)
	worker := s.workers[uid]
	if worker == nil {
		return fmt.Errorf("no worker exists of target:%s, pattern:%s",
			target, pattern)
	}
	worker.stop()
	delete(s.workers, uid)
	return nil
}

// GetResult returns the building result
func (s *serverImpl) GetResult(target builder.Target,
	pattern builder.Pattern,
	from, to uint,
) (interface{}, error) {
	uid := workerUID(target, pattern)
	if s.workers[uid] == nil {
		return nil, fmt.Errorf("no worker exists of target:%s, pattern:%s",
			target, pattern)
	}

	codec := mq.GetCodec(target, pattern)
	if codec == nil {
		return nil, fmt.Errorf("no codec of target:%s,pattern:%s", target, pattern)
	}

	msgs, err := s.mq.Get(uid, from, to, codec)
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, len(msgs))
	for i, m := range msgs {
		ret[i] = m.Body
	}

	return ret, nil
}

// Close workers
func (s *serverImpl) Close() error {
	for uid, w := range s.workers {
		_ = uid
		w.stop()
		delete(s.workers, uid)
	}
	return nil
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
		impl:    impl,
		mq:      q,
		workers: make(map[string]worker),
		logger:  logger.Std,
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
	return srv, nil
}
