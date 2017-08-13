package worker

import (
	"qiniu.com/video/builder"
	"qiniu.com/video/mq"
)

// Server interface
type Server interface {
	StartBuild(target builder.Target, pattern builder.Pattern, mq mq.MQ) error
	StopBuild(target builder.Target, pattern builder.Pattern) error
	Close() error
}

type workerImpl struct {
	impl builder.Implement
}

// StartBuild the building
func (w *workerImpl) StartBuild(target builder.Target,
	pattern builder.Pattern,
	mq mq.MQ,
) error {
	// TODO
	return nil
}

// StopBuild the building
func (w *workerImpl) StopBuild(target builder.Target,
	pattern builder.Pattern,
) error {
	// TODO
	return nil
}

// Close workers
func (w *workerImpl) Close() error {
	// TODO
	return nil
}

// CreateWorker create build worker
func CreateWorker(impl builder.Implement) (Server, error) {
	// TODO
	return nil
}

type worker struct {
	close   <-chan int
	target  builder.Target
	pattern builder.Pattern
	mq      mq.MQ
}

func (w *worker) start() {
	// TODO
}

func (w *worker) stop() {
	w.close <- 0
}
