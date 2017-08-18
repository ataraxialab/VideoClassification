package mq

import (
	"encoding/binary"
	"time"

	"qiniu.com/video/config"

	"github.com/boltdb/bolt"
)

const (
	dataPath = "./embeddedmq.db"
)

// EmbeddedMQ embedded mq using `bolt`
type EmbeddedMQ struct {
	db     *bolt.DB
	endian binary.ByteOrder
}

// Open mq
func (mq *EmbeddedMQ) Open(conf *config.MQ) error {
	db, err := bolt.Open(conf.URI, 0600, nil)
	if err != nil {
		return err
	}
	mq.db = db
	mq.endian = binary.BigEndian
	return nil
}

// Close mq
func (mq *EmbeddedMQ) Close() error {
	return mq.db.Close()
}

// Put messages
func (mq *EmbeddedMQ) Put(topic string,
	encoder Encoder,
	val ...interface{}) error {
	db := mq.db
	err := db.Batch(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(topic))
		if err != nil {
			return err
		}

		k, _ := bucket.Cursor().Last()
		id := uint64(0)
		if k != nil {
			id = mq.endian.Uint64(k)
			id++
		}

		createdAt := time.Now().Unix()
		for _, v := range val {
			m := MessageEx{
				ID:        mq.messageID(id),
				Status:    StatusPending,
				CreatedAt: uint64(createdAt),
				Body:      v,
			}
			err = bucket.Put(mq.messageID(id), m.Encode(encoder))
			if err != nil {
				return err
			}
			id++
		}
		return nil
	})
	return err
}

// Get messages
func (mq *EmbeddedMQ) Get(topic string, from, count uint, decoder Decoder) (
	[]MessageEx, error) {

	db := mq.db
	var vals []MessageEx

	err := db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(topic))
		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()
		min := from
		max := min + count
		k := mq.messageID(uint64(from))
		vals = make([]MessageEx, 0, count)
		for k, v := c.Seek(k); k != nil && min < max; k, v = c.Next() {
			min++
			m := MessageEx{}
			m.Decode(v, decoder)
			vals = append(vals, m)
		}

		return nil
	})
	return vals, err
}

// Delete message
func (mq *EmbeddedMQ) Delete(topic string, ids ...[]byte) error {
	db := mq.db
	err := db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(topic))
		if bucket == nil {
			return nil
		}

		for _, k := range ids {
			if err := bucket.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// DeleteTopic deletes topic
func (mq *EmbeddedMQ) DeleteTopic(topic string) error {
	db := mq.db
	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(topic))
	})
}

func (mq *EmbeddedMQ) messageID(id uint64) []byte {
	bytes := make([]byte, 8)
	mq.endian.PutUint64(bytes, id)
	return bytes
}
