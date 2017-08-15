package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"qiniu.com/video/builder"
	"qiniu.com/video/mq"
)

type mockWorker int

func (mw *mockWorker) start() {
	*mw++
}

func (mw *mockWorker) stop() {
	*mw--
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

type mockBuilder struct{}

func (mockBuilder) Build(string, interface{}) ([]interface{}, error) {
	return nil, nil
}

const (
	impl    = builder.Implement("mockImpl")
	target  = builder.Target("mockTarget")
	pattern = builder.Pattern("mockPattern")
)

func init() {
	builder.Register(impl, target, pattern, mockBuilder{})
	mq.Register(target, pattern, mockCodec{})
}

func TestServer(t *testing.T) {
	server := &mockServer{
		serverImpl: &serverImpl{
			impl:         impl,
			workers:      make(map[string]worker),
			createWorker: createWorker,
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

type mockMQ struct{}

func (mq mockMQ) Open() error {
	return nil
}

func (mq mockMQ) Close() error {
	return nil
}

func (mq mockMQ) Put(topic string,
	encoder mq.Encoder,
	val ...interface{},
) error {
	return nil
}

func (mq mockMQ) Get(topic string, from, count uint,
	decoder mq.Decoder) ([]mq.MessageEx, error) {
	return nil, nil
}

func (mq mockMQ) Delete(topic string, ids ...[]byte) error {
	return nil
}

func TestCreateServer(t *testing.T) {
	_, err := CreateServer(builder.Implement(""), nil)
	assert.NotNil(t, err)
	_, err = CreateServer(impl, nil)
	assert.NotNil(t, err)
	srv, err := CreateServer(impl, mockMQ{})
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	var _ Server = srv
}
