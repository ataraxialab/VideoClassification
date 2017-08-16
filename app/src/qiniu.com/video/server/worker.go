package server

import (
	"sync/atomic"
	"time"

	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/mq"
)

type worker interface {
	start()
	stop()
	pause()
	proceed()
	status() workerStatus
	name() string
}

type workerStatus interface {
	buildCount() uint64
	status() uint8
}

type vWorkerStatus interface {
	workerStatus
	cmsStatus(old, new uint8) bool
	addCount(delta int) uint64
}

const (
	// Start starting status
	Start uint8 = iota + 1
	// Pause pause status
	Pause
	// Stop closed status
	Stop
)

type statusImpl struct {
	count uint64
	stat  uint32
}

func (s *statusImpl) buildCount() uint64 {
	return atomic.LoadUint64(&s.count)
}

func (s *statusImpl) status() uint8 {
	return uint8(atomic.LoadUint32(&s.stat))
}

func (s *statusImpl) cmsStatus(old, new uint8) bool {
	return atomic.CompareAndSwapUint32(&s.stat, uint32(old), uint32(new))
}
func (s *statusImpl) addCount(delta int) uint64 {
	return atomic.AddUint64(&s.count, uint64(delta))
}

type workerImpl struct {
	uid         string
	mq          mq.MQ
	codec       mq.Codec
	params      interface{}
	dataBuilder builder.Builder
	logger      *logger.Logger
	goon        chan int
	stat        vWorkerStatus
}

func (w *workerImpl) start() {
	if w.stat != nil && !w.stat.cmsStatus(uint8(0), Start) {
		w.logger.Errorf("invalid start operation, cur status:%d",
			w.stat.status())
		return
	}

	w.logger.Infof("%s start working", w.uid)
	w.goon = make(chan int, 1)
	w.stat = &statusImpl{stat: uint32(Start)}

	go func() {
		w.logger.Info("start clean")
		from, count, expireTime := uint(0), uint(100), int64(time.Minute*30)
		for w.stat.status() != Stop {
			time.Sleep(5 * time.Minute)
			w.clean(from, count, expireTime)
		}
	}()

	go func() {
		for w.stat.status() != Stop {
			if w.stat.status() == Pause {
				<-w.goon
			}

			ret, err := w.dataBuilder.Build(w.params)
			if err != nil {
				logger.Errorf("build error:%v", err)
				continue
			}

			w.stat.addCount(len(ret))

			err = w.mq.Put(w.uid, w.codec, ret...)
			if err != nil {
				logger.Errorf("put message to mq error:%v", err)
			}
		}
		close(w.goon)
	}()
}

func (w *workerImpl) pause() {
	logger.Infof("pause:%s, %d", w.uid, w.stat.status())
	if !w.stat.cmsStatus(Start, Pause) {
		logger.Infof("invalid pause, old:%d", w.stat.status())
	}
}

func (w *workerImpl) stop() {
	logger.Infof("stop:%s, %d", w.uid, w.stat.status())
	if !w.stat.cmsStatus(Start, Stop) && !w.stat.cmsStatus(Pause, Stop) {
		logger.Infof("invalid stop, old:%d", w.stat.status())
	}
}

func (w *workerImpl) proceed() {
	logger.Infof("proceed:%s, %d", w.uid, w.stat.status())
	if !w.stat.cmsStatus(Pause, Start) {
		logger.Infof("invalid proceed, old:%d", w.stat.status())
		return
	}
	w.goon <- 0
}

func (w *workerImpl) clean(from, count uint, expireTime int64) {
	msgs, err := w.mq.Get(w.uid, from, count, w.codec)
	if err != nil {
		w.logger.Errorf("clean get data error:%v", err)
		return
	}

	if len(msgs) == 0 {
		return
	}

	keys := make([][]byte, 0, len(msgs))
	now := time.Now().Unix()
	for _, m := range msgs {
		if now-int64(m.CreatedAt) < expireTime {
			break
		}
		err := w.dataBuilder.Clean(m.Body)
		if err != nil {
			w.logger.Errorf("clean resource error:%v", m.Body)
			continue
		}
		keys = append(keys, m.ID)
	}

	w.mq.Delete(w.uid, keys...)
}

func (w *workerImpl) status() workerStatus {
	return w.stat
}

func (w *workerImpl) name() string {
	return w.uid
}
