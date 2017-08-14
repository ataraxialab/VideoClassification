package server

import (
	"qiniu.com/video/builder"
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
}

func (w *workerImpl) start() {
	go func() {
		for !w.isClosed {
			// TODO
		}
	}()
}

func (w *workerImpl) stop() {
	w.isClosed = true
}
