package mq

import (
	"encoding/binary"
	"time"

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
func (mq *EmbeddedMQ) Open() error {
	db, err := bolt.Open(dataPath, 0600, nil)
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
func (mq *EmbeddedMQ) Put(topic string, val ...Message) error {
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
			m := message{
				id:        mq.messageID(id),
				status:    StatusPending,
				createdAt: uint64(createdAt),
				body:      v.Encode(),
			}
			err = bucket.Put(mq.messageID(id), m.Encode())
			if err != nil {
				return err
			}
			id++
		}
		return nil
	})
	return err
}

func (mq *EmbeddedMQ) getMessage(topic string, from, count uint) ([][]byte, error) {
	db := mq.db
	var vals [][]byte

	err := db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(topic))
		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()
		min := from
		max := min + count
		k := mq.messageID(uint64(from))
		vals = make([][]byte, 0, count)
		for k, v := c.Seek(k); k != nil && min < max; k, v = c.Next() {
			min++
			vals = append(vals, v)
		}

		return nil
	})
	return vals, err
}

// Get messages
func (mq *EmbeddedMQ) Get(topic string, from, count uint) ([][]byte, error) {
	rawMsgs, err := mq.getMessage(topic, from, count)
	if err != nil {
		return nil, err
	}

	vals := make([][]byte, 0, len(rawMsgs))
	m := message{}
	for _, rm := range rawMsgs {
		m.Decode(rm)
		vals = append(vals, m.body)
	}

	return vals, nil
}

func (mq *EmbeddedMQ) messageID(id uint64) []byte {
	bytes := make([]byte, 8)
	mq.endian.PutUint64(bytes, id)
	return bytes
}
