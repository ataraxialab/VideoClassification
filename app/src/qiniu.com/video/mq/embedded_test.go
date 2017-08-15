package mq

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stringCodec struct{}

func (sc stringCodec) Encode(s interface{}) []byte {
	return []byte(s.(string))
}

func (sc stringCodec) Decode(bytes []byte) interface{} {
	return string(bytes)
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
	msgs := []string{"1", "2", "3", "4", "5", "6"}
	wmsgs := make([]interface{}, len(msgs))
	for i, m := range msgs {
		wmsgs[i] = m
	}

	codec := stringCodec{}
	assert.Nil(t, mq.Put(topic, codec, wmsgs...))

	t.Run("get", func(t *testing.T) {
		rawMsgs, err := mq.Get(topic, 0, 3, codec)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(rawMsgs))

		for i, m := range rawMsgs {
			assert.Equal(t, msgs[i], m.Body)
			assert.Equal(t, mq.messageID(uint64(i)), m.ID)
		}
		rawMsgs, err = mq.Get(topic, 1, 8, codec)
		assert.Nil(t, err)
		assert.Equal(t, 5, len(rawMsgs))

		for i, m := range rawMsgs {
			assert.Equal(t, msgs[i+1], m.Body)
			assert.Equal(t, mq.messageID(uint64(i+1)), m.ID)
		}
	})

	t.Run("delete", func(t *testing.T) {
		err := mq.Delete(topic, mq.messageID(uint64(4)), mq.messageID(uint64(5)))
		assert.Nil(t, err)
		rawMsgs, err := mq.Get(topic, 1, 8, codec)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(rawMsgs))
	})

	t.Run("delete-topic", func(t *testing.T) {
		err := mq.DeleteTopic(topic)
		assert.Nil(t, err)
		rawMsgs, err := mq.Get(topic, 1, 8, codec)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(rawMsgs))
		err = mq.DeleteTopic(topic)
		assert.NotNil(t, err)
	})
}
