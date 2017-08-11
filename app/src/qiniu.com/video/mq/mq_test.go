package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ByteCodec []byte

func (bc ByteCodec) Encode(bytes interface{}) []byte {
	return bytes.([]byte)
}

func (bc ByteCodec) Decode(bytes []byte) interface{} {
	return bytes
}

var codec = ByteCodec{}

func TestMessageEx(t *testing.T) {
	m := MessageEx{
		ID:        []byte("id"),
		CreatedAt: 9999,
		Status:    StatusDeleted,
		Body:      []byte("Body"),
	}

	md := MessageEx{}
	md.Decode(m.Encode(codec), codec)

	assert.Equal(t, m.ID, md.ID)
	assert.Equal(t, m.CreatedAt, md.CreatedAt)
	assert.Equal(t, uint16(StatusDeleted), md.Status)
	assert.Equal(t, m.Body, md.Body)
}

func TestUpdateStatus(t *testing.T) {
	m := MessageEx{
		ID:        []byte("id"),
		CreatedAt: 9999,
		Status:    StatusDeleted,
		Body:      []byte("Body"),
	}

	bytes := m.Encode(codec)
	updateStatus(bytes, StatusPending)

	m.Decode(bytes, codec)
	assert.Equal(t, uint16(StatusPending), m.Status)
}
