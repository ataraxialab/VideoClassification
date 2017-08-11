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
	Get(topic string, from, count uint) ([][]byte, error)
}

// Message represent the value stored in the MQ
type Message interface {
	Encode() []byte
	Decode([]byte)
}

// BaseMessage message common data
type message struct {
	createdAt uint64
	status    uint16
	id        []byte
	body      []byte
}

var endian = binary.BigEndian

// Encode the base message to bytes
func (m *message) Encode() []byte {
	idLen, bodyLen := len(m.id), len(m.body)
	bytes := make([]byte, 2+8+2+idLen+bodyLen)
	endian.PutUint16(bytes, m.status)
	endian.PutUint64(bytes[2:], m.createdAt)
	endian.PutUint16(bytes[2+8:], uint16(idLen))
	copy(bytes[2+8+2:], m.id)
	copy(bytes[2+8+2+idLen:], m.body)

	return bytes
}

// Decode the raw bytes to `BaseMessage`
func (m *message) Decode(bytes []byte) {
	m.status = endian.Uint16(bytes)
	m.createdAt = endian.Uint64(bytes[2:])
	idLen := endian.Uint16(bytes[2+8:])
	m.id = bytes[2+8+2 : 2+8+2+idLen]
	m.body = bytes[2+8+2+idLen:]
}
