package server

import (
	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/mq"
)

type worker interface {
	start()
	stop()
	pause()
	proceed()
}

const (
	// Start starting status
	Start = iota
	// Pause pause status
	Pause
	// Stop closed status
	Stop
)

type workerImpl struct {
	uid         string
	mq          mq.MQ
	codec       mq.Codec
	params      interface{}
	dataBuilder builder.Builder
	logger      *logger.Logger
	status      int
	goon        chan int
}

func (w *workerImpl) start() {
	w.goon = make(chan int, 1)
	go func() {
		for w.status != Stop {
			if w.status == Pause {
				<-w.goon
			}

			ret, err := w.dataBuilder.Build(w.selectVideo(), w.params)
			if err != nil {
				logger.Errorf("build error:%v", err)
				continue
			}

			err = w.mq.Put(w.uid, w.codec, ret...)
			if err != nil {
				logger.Errorf("put message to mq error:%v", err)
			}
		}
	}()
}

func (w *workerImpl) pause() {
	logger.Infof("pause:%s", w.uid)
	w.status = Pause
}

func (w *workerImpl) stop() {
	logger.Infof("stop:%s", w.uid)
	w.status = Stop
}
func (w *workerImpl) proceed() {
	w.status = Start
	w.goon <- 0
}

func (w *workerImpl) selectVideo() string {
	// TODO
	return ""
}
