package server

import (
	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/mq"
)

type worker interface {
	start()
	stop()
}

type workerImpl struct {
	uid         string
	mq          mq.MQ
	codec       mq.Codec
	params      interface{}
	dataBuilder builder.Builder
	isClosed    bool
	logger      *logger.Logger
}

func (w *workerImpl) start() {
	go func() {
		for !w.isClosed {
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

func (w *workerImpl) stop() {
	logger.Infof("stop:%s", w.uid)
	w.isClosed = true
}

func (w *workerImpl) selectVideo() string {
	// TODO
	return ""
}
