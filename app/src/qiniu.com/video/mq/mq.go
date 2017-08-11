package mq

import "encoding/binary"

const (
	// StatusPending waiting to consume
	StatusPending = iota
	// StatusConsuming is consuming
	StatusConsuming
	// StatusDeleted is deleted
	StatusDeleted
)

// MQ provides the persistent operation
type MQ interface {
	Open() error
	Close() error
	Put(topic string, val ...Message) error
	Get(topic string, from, count int, messages []Message) (int, error)
}

// Message represent the value stored in the MQ
type Message interface {
	Encode() []byte
	Decode([]byte)

	SetID(id string)
	ID() string
	SetCreatedAt(uint64)
	CreatedAt() uint64
	SetStatus(int)
	Status()
}

// BaseMessage message common data
type BaseMessage struct {
	id        string
	createdAt uint64
	status    uint16
}

// SetID store the message id
func (bm *BaseMessage) SetID(id string) {
	bm.id = id
}

// ID return message id
func (bm *BaseMessage) ID() string {
	return bm.id
}

// SetCreatedAt set create time
func (bm *BaseMessage) SetCreatedAt(time uint64) {
	bm.createdAt = time
}

// CreatedAt return message created time
func (bm *BaseMessage) CreatedAt() uint64 {
	return bm.createdAt
}

// SetStatus set message status
func (bm *BaseMessage) SetStatus(status uint16) {
	bm.status = status
}

// Status return message status
func (bm *BaseMessage) Status() uint16 {
	return bm.status
}

var endian = binary.BigEndian

// Encode the base message to bytes
func (bm *BaseMessage) Encode() []byte {
	idLen := len(bm.id)
	bytes := make([]byte, 2+idLen+8+2)
	endian.PutUint16(bytes, uint16(idLen))
	copy(bytes[2:], bm.id)
	endian.PutUint64(bytes[2+idLen:], bm.createdAt)
	endian.PutUint16(bytes[2+idLen+8:], bm.status)

	return bytes
}

// Decode the raw bytes to `BaseMessage`
func (bm *BaseMessage) Decode(bytes []byte) {
	idLen := endian.Uint16(bytes)
	bm.id = string(bytes[2 : 2+idLen])
	bm.createdAt = endian.Uint64(bytes[2+idLen:])
	bm.status = endian.Uint16(bytes[2+idLen+8:])
}
