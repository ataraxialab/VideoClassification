package server

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/mq"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

type mockWorker struct {
	counter int
	uid     string
}

var status = &statusImpl{}

func (mw *mockWorker) start() {
	mw.counter++
	go func() {
		status.addCount(1)
	}()
}

func (mw *mockWorker) stop() {
	mw.counter--
}

func (mw *mockWorker) pause() {
}

func (mw *mockWorker) proceed() {
}

func (mw *mockWorker) name() string {
	return mw.uid
}

func (mw *mockWorker) status() workerStatus {
	return status
}

func createWorker(name string, p interface{}, b builder.Builder,
	q mq.Codec) worker {
	w := &mockWorker{}
	w.uid = name
	return w
}

type mockCodec struct{}

func (mockCodec) Encode(interface{}) []byte {
	return nil
}
func (mockCodec) Decode([]byte) interface{} {
	return nil
}

type mockBuilder int

func (b *mockBuilder) Build(interface{}) ([]interface{}, error) {
	*b++
	return make([]interface{}, 100000), nil
}

func (b *mockBuilder) Clean(interface{}) error {
	*b++
	return nil
}

const (
	impl  = builder.Implement("mockImpl")
	targt = target.Target("mockTarget")
	pat   = pattern.Pattern("mockPattern")
)

func init() {
	b := mockBuilder(0)
	builder.Register(impl, targt, pat, &b)
	mq.Register(targt, pat, mockCodec{})
}

func TestStopStart(t *testing.T) {
	server := &serverImpl{
		impl:           impl,
		workers:        make(map[string]worker),
		consumedCounts: make(map[string]uint64),
		createWorker:   createWorker,
		mq:             &mockMQ{},
		logger:         logger.Std,
		locker:         new(sync.Mutex),
	}
	err := server.StartBuilding(targt, pat, nil)
	assert.Nil(t, err)
	worker := server.workers[workerUID(targt, pat)]
	if worker == nil {
		t.Fatal("nil worker")
	}
	assert.Equal(t, 1, (worker.(*mockWorker).counter))
	err = server.StartBuilding(targt, pat, nil)
	assert.NotNil(t, err)

	err = server.StopBuilding(targt, pat)
	assert.Nil(t, err)

	err = server.StopBuilding(targt, pat)
	assert.NotNil(t, err)

	err = server.StartBuilding(targt, pat, nil)
	assert.Nil(t, err)
	err = server.StopBuilding(targt, pat)
	assert.Nil(t, err)

	err = server.StartBuilding(targt, pat, nil)
	assert.Nil(t, err)
	server.Close()
	assert.Equal(t, 0, len(server.workers))
}

type mockMQ struct {
	deleteCount int
}

func (q *mockMQ) Open() error {
	return nil
}

func (q *mockMQ) Close() error {
	return nil
}

func (q *mockMQ) Put(topic string,
	encoder mq.Encoder,
	val ...interface{},
) error {
	return nil
}

func (q *mockMQ) Get(topic string, from, count uint,
	decoder mq.Decoder) ([]mq.MessageEx, error) {
	return []mq.MessageEx{
		mq.MessageEx{
			Body:      "test",
			ID:        []byte("1"),
			CreatedAt: uint64(time.Now().Unix() - int64(time.Millisecond*200)),
		},
		mq.MessageEx{
			Body:      "test",
			ID:        []byte("1"),
			CreatedAt: uint64(time.Now().Unix() - int64(time.Millisecond*100)),
		},
		mq.MessageEx{
			Body:      "test",
			ID:        []byte("1"),
			CreatedAt: uint64(time.Now().Unix() + int64(time.Hour)),
		},
	}, nil
}

func (q *mockMQ) Delete(topic string, ids ...[]byte) error {
	q.deleteCount = len(ids)
	return nil
}

func (q mockMQ) DeleteTopic(topic string) error {
	return nil
}

func TestCreateServer(t *testing.T) {
	_, err := CreateServer(builder.Implement(""), nil)
	assert.NotNil(t, err)
	_, err = CreateServer(impl, nil)
	assert.NotNil(t, err)
	srv, err := CreateServer(impl, &mockMQ{})
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	var _ Server = srv
}

func TestGetResult(t *testing.T) {
	srv, err := CreateServer(impl, &mockMQ{})
	var _ Server = srv
	assert.Nil(t, err)

	_, err = srv.GetResult(targt, pat, 0, 1)
	assert.NotNil(t, err)

	err = srv.StartBuilding(targt, pat, nil)
	assert.Nil(t, err)

	ret, err := srv.GetResult(targt, pat, 0, 1)
	slice, ok := ret.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, slice[0], "test")

	srv.Close()
}

func TestWorker(t *testing.T) {
	builder := mockBuilder(0)
	w := &workerImpl{
		uid:         "test-worker",
		mq:          &mockMQ{},
		codec:       mockCodec{},
		dataBuilder: &builder,
		logger:      logger.Std,
	}

	w.start()
	time.Sleep(100 * time.Millisecond)
	assert.True(t, builder > 0)
	assert.Equal(t, w.stat.status(), Start)

	w.start()
	assert.Equal(t, w.stat.status(), Start)

	w.proceed()
	assert.Equal(t, w.stat.status(), Start)

	w.pause()
	cnt := builder
	time.Sleep(100 * time.Millisecond)
	assert.True(t, cnt == builder || cnt+1 == builder)
	assert.Equal(t, w.stat.status(), Pause)

	w.start()
	assert.Equal(t, w.stat.status(), Pause)

	w.proceed()
	time.Sleep(1000 * time.Millisecond)
	assert.True(t, builder > cnt+1)
	assert.Equal(t, w.stat.status(), Start)

	w.stop()
	cnt = builder
	time.Sleep(1000 * time.Millisecond)
	assert.True(t, cnt == builder || cnt+1 == builder)
	assert.Equal(t, w.stat.status(), Stop)

	w.start()
	assert.Equal(t, w.stat.status(), Stop)
	w.proceed()
	assert.Equal(t, w.stat.status(), Stop)
}

func TestClean(t *testing.T) {
	builder := mockBuilder(0)
	q := &mockMQ{}
	w := &workerImpl{
		uid:         "test-worker",
		mq:          q,
		codec:       mockCodec{},
		dataBuilder: &builder,
		logger:      logger.Std,
	}

	builder = 10
	w.clean(uint(0), uint(3), int64(400*time.Millisecond))
	assert.Equal(t, 10, int(builder))
	assert.Equal(t, 0, q.deleteCount)

	builder = 0
	w.clean(uint(0), uint(3), int64(200*time.Millisecond))
	assert.Equal(t, 1, int(builder))
	assert.Equal(t, 1, q.deleteCount)

	builder = 0
	w.clean(uint(0), uint(3), int64(50*time.Millisecond))
	assert.Equal(t, 2, int(builder))
	assert.Equal(t, 2, q.deleteCount)
}

func TestMonitor(t *testing.T) {
	server := &serverImpl{
		impl:           impl,
		workers:        make(map[string]worker),
		consumedCounts: make(map[string]uint64),
		mq:             &mockMQ{},
		locker:         new(sync.Mutex),
		maxRetainCount: 100,
		monitorPeriod:  10 * time.Millisecond,
		logger:         logger.Std,
	}
	server.createWorker = func(uid string, params interface{},
		dataBuilder builder.Builder, codec mq.Codec) worker {
		return &workerImpl{
			uid:         uid,
			mq:          server.mq,
			codec:       codec,
			params:      params,
			dataBuilder: dataBuilder,
			logger:      server.logger,
		}
	}
	go server.monitor()
	var _ Server = server

	err := server.StartBuilding(targt, pat, nil)
	assert.Nil(t, err)
	uid := workerUID(targt, pat)
	worker := server.workers[uid]
	if worker == nil {
		t.Fatal("nil worker")
	}
	time.Sleep(500 * time.Millisecond)
	s := worker.status()
	t.Logf("build count:%d, consume:%d", s.buildCount(), server.consumedCounts[uid])
	assert.Equal(t, Pause, s.status())
	for s.status() != Start {
		server.GetResult(targt, pat, 0, 100)
	}
	assert.Equal(t, Start, s.status())
	t.Logf("build count:%d, consume:%d", s.buildCount(), server.consumedCounts[uid])
	time.Sleep(500 * time.Millisecond)
	t.Logf("build count:%d, consume:%d", s.buildCount(), server.consumedCounts[uid])
	assert.Equal(t, Pause, s.status())
	t.Logf("build count:%d, consume:%d", s.buildCount(), server.consumedCounts[uid])

	server.Close()
}
