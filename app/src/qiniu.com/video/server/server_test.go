package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"qiniu.com/video/builder"
	"qiniu.com/video/logger"
	"qiniu.com/video/mq"
)

type mockWorker int

func (mw *mockWorker) start() {
	*mw++
}

func (mw *mockWorker) stop() {
	*mw--
}

func (mw *mockWorker) pause() {
}

func (mw *mockWorker) proceed() {
}

type mockServer struct {
	*serverImpl
}

func createWorker(string, interface{}, builder.Builder, mq.Codec) worker {
	w := mockWorker(0)
	return &w
}

type mockCodec struct{}

func (mockCodec) Encode(interface{}) []byte {
	return nil
}
func (mockCodec) Decode([]byte) interface{} {
	return nil
}

type mockBuilder int

func (b *mockBuilder) Build(string, interface{}) ([]interface{}, error) {
	*b++
	return nil, nil
}

func (b *mockBuilder) Clean(interface{}) error {
	*b++
	return nil
}

const (
	impl    = builder.Implement("mockImpl")
	target  = builder.Target("mockTarget")
	pattern = builder.Pattern("mockPattern")
)

func init() {
	b := mockBuilder(0)
	builder.Register(impl, target, pattern, &b)
	mq.Register(target, pattern, mockCodec{})
}

func TestServer(t *testing.T) {
	server := &mockServer{
		serverImpl: &serverImpl{
			impl:         impl,
			workers:      make(map[string]worker),
			createWorker: createWorker,
			mq:           &mockMQ{},
		},
	}
	err := server.StartBuilding(target, pattern, nil)
	assert.Nil(t, err)
	worker := server.workers[workerUID(target, pattern)]
	if worker == nil {
		t.Fatal("nil worker")
	}
	assert.Equal(t, 1, int(*(worker.(*mockWorker))))
	err = server.StartBuilding(target, pattern, nil)
	assert.NotNil(t, err)

	err = server.StopBuilding(target, pattern)
	assert.Nil(t, err)

	err = server.StopBuilding(target, pattern)
	assert.NotNil(t, err)

	err = server.StartBuilding(target, pattern, nil)
	assert.Nil(t, err)
	err = server.StopBuilding(target, pattern)
	assert.Nil(t, err)

	err = server.StartBuilding(target, pattern, nil)
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

	_, err = srv.GetResult(target, pattern, 0, 1)
	assert.NotNil(t, err)

	err = srv.StartBuilding(target, pattern, nil)
	assert.Nil(t, err)

	ret, err := srv.GetResult(target, pattern, 0, 1)
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
	assert.Equal(t, w.status, Start)
	w.pause()
	cnt := builder
	time.Sleep(100 * time.Millisecond)
	assert.True(t, cnt == builder || cnt+1 == builder)
	assert.Equal(t, w.status, Pause)
	w.proceed()
	time.Sleep(1000 * time.Millisecond)
	assert.True(t, builder > cnt+1)
	assert.Equal(t, w.status, Start)
	w.stop()
	cnt = builder
	time.Sleep(1000 * time.Millisecond)
	assert.True(t, cnt == builder || cnt+1 == builder)
	assert.Equal(t, w.status, Stop)
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
