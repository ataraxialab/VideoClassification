package server

import (
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
		w.logger.Info("start clean")
		from, count, expireTime := uint(0), uint(100), int64(time.Minute*30)
		for w.status != Stop {
			time.Sleep(5 * time.Minute)
			w.clean(from, count, expireTime)
		}
	}()
	go func() {
		for w.status != Stop {
			if w.status == Pause {
				<-w.goon
			}

			ret, err := w.dataBuilder.Build(w.params)
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
