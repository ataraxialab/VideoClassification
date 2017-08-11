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
	Get(topic string, from, count uint) ([]MessageEx, error)
	Delete(topic string, ids ...[]byte) error
}

// Message represent the value stored in the MQ
type Message interface {
	Encode() []byte
	Decode([]byte)
}

// MessageEx message common data
type MessageEx struct {
	CreatedAt uint64
	Status    uint16
	ID        []byte
	Body      []byte
}

var endian = binary.BigEndian

// Encode the base message to bytes
func (m *MessageEx) Encode() []byte {
	idLen := len(m.ID)
	bytes := make([]byte, 2+8+2+idLen+len(m.Body))
	endian.PutUint16(bytes, m.Status)
	endian.PutUint64(bytes[2:], m.CreatedAt)
	endian.PutUint16(bytes[2+8:], uint16(idLen))
	copy(bytes[2+8+2:], m.ID)
	copy(bytes[2+8+2+idLen:], m.Body)

	return bytes
}

// Decode the raw bytes to `BaseMessage`
func (m *MessageEx) Decode(bytes []byte) {
	m.Status = endian.Uint16(bytes)
	m.CreatedAt = endian.Uint64(bytes[2:])
	idLen := endian.Uint16(bytes[2+8:])
	m.ID = bytes[2+8+2 : 2+8+2+idLen]
	m.Body = bytes[2+8+2+idLen:]
}

func updateStatus(rawMessage []byte, status uint16) []byte {
	endian.PutUint16(rawMessage, status)
	return rawMessage
}
