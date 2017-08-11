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
	mq := &EmbeddedMQ{}
	var _ MQ = mq

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

	t.Run("get", func(t *testing.T) {
		rawMsgs, err := mq.Get(topic, 0, 3)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(rawMsgs))

		for i, m := range rawMsgs {
			assert.Equal(t, []byte(msgs[i]), m.Body)
			assert.Equal(t, mq.messageID(uint64(i)), m.ID)
		}
		rawMsgs, err = mq.Get(topic, 1, 8)
		assert.Nil(t, err)
		assert.Equal(t, 5, len(rawMsgs))

		for i, m := range rawMsgs {
			assert.Equal(t, []byte(msgs[i+1]), m.Body)
			assert.Equal(t, mq.messageID(uint64(i+1)), m.ID)
		}
	})

	t.Run("delete", func(t *testing.T) {
		err := mq.Delete(topic, mq.messageID(uint64(4)), mq.messageID(uint64(5)))
		assert.Nil(t, err)
		rawMsgs, err := mq.Get(topic, 1, 8)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(rawMsgs))
	})
}
