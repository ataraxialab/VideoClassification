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
	Put(topic string, encoder Encoder, val ...interface{}) error
	Get(topic string, from, count uint, decoder Decoder) ([]MessageEx, error)
	Delete(topic string, ids ...[]byte) error
}

// Decoder decode the bytes to user type
type Decoder interface {
	Decode([]byte) interface{}
}

// Encoder encode message to byte
type Encoder interface {
	Encode(interface{}) []byte
}

// MessageEx message common data
type MessageEx struct {
	CreatedAt uint64
	Status    uint16
	ID        []byte
	Body      interface{}
}

var endian = binary.BigEndian

// Encode the base message to bytes
func (m *MessageEx) Encode(encoder Encoder) []byte {
	idLen := len(m.ID)
	body := encoder.Encode(m.Body)
	bytes := make([]byte, 2+8+2+idLen+len(body))
	endian.PutUint16(bytes, m.Status)
	endian.PutUint64(bytes[2:], m.CreatedAt)
	endian.PutUint16(bytes[2+8:], uint16(idLen))
	copy(bytes[2+8+2:], m.ID)
	copy(bytes[2+8+2+idLen:], body)

	return bytes
}

// Decode the raw bytes to `BaseMessage`
func (m *MessageEx) Decode(bytes []byte, decoder Decoder) {
	m.Status = endian.Uint16(bytes)
	m.CreatedAt = endian.Uint64(bytes[2:])
	idLen := endian.Uint16(bytes[2+8:])
	m.ID = bytes[2+8+2 : 2+8+2+idLen]
	m.Body = decoder.Decode(bytes[2+8+2+idLen:])
}

func updateStatus(rawMessage []byte, status uint16) []byte {
	endian.PutUint16(rawMessage, status)
	return rawMessage
}
