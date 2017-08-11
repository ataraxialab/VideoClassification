package mq

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMsg string

func (m mockMsg) Encode() []byte {
	return []byte(m)
}

func (m mockMsg) Decode(bytes []byte) {
}

func TestEmbedded(t *testing.T) {
	var mq MQ
	mq = &EmbeddedMQ{}
	mq.Open()
	defer func() {
		mq.Close()
		os.Remove(dataPath)
	}()

	topic := "topic"
	msgs := []mockMsg{"1", "2", "3", "4", "5", "6"}
	wmsgs := make([]Message, len(msgs))
	for i, m := range msgs {
		wmsgs[i] = m
	}

	assert.Nil(t, mq.Put(topic, wmsgs...))

	rawMsgs, err := mq.Get(topic, 0, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(rawMsgs))

	for i, m := range rawMsgs {
		assert.Equal(t, []byte(msgs[i]), m)
	}
	rawMsgs, err = mq.Get(topic, 1, 8)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(rawMsgs))

	for i, m := range rawMsgs {
		assert.Equal(t, []byte(msgs[i+1]), m)
	}
}
