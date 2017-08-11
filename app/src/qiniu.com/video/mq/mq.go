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
	Get(topic string, from, count uint, messages *[]Message) (uint, error)
}

// Message represent the value stored in the MQ
type Message interface {
	Encode() []byte
	Decode([]byte)
}

// BaseMessage message common data
type message struct {
	id        []byte
	createdAt uint64
	status    uint16
	body      []byte
}

var endian = binary.BigEndian

// Encode the base message to bytes
func (m *message) Encode() []byte {
	idLen, bodyLen := len(m.id), len(m.body)
	bytes := make([]byte, 2+idLen+8+2+bodyLen)
	endian.PutUint16(bytes, uint16(idLen))
	copy(bytes[2:], m.id)
	endian.PutUint64(bytes[2+idLen:], m.createdAt)
	endian.PutUint16(bytes[2+idLen+8:], m.status)
	copy(bytes[2+idLen+8+2:], m.body)

	return bytes
}

// Decode the raw bytes to `BaseMessage`
func (m *message) Decode(bytes []byte) {
	idLen := endian.Uint16(bytes)
	m.id = bytes[2 : 2+idLen]
	m.createdAt = endian.Uint64(bytes[2+idLen:])
	m.status = endian.Uint16(bytes[2+idLen+8:])
	m.body = bytes[2+idLen+8+2:]
}
