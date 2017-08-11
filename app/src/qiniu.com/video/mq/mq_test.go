package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseMessage(t *testing.T) {
	m := message{
		id:        []byte("id"),
		createdAt: 9999,
		status:    StatusDeleted,
		body:      []byte("body"),
	}

	md := message{}

	md.Decode(m.Encode())

	assert.Equal(t, m.id, md.id)
	assert.Equal(t, m.createdAt, md.createdAt)
	assert.Equal(t, uint16(StatusDeleted), md.status)
	assert.Equal(t, m.body, md.body)
}

func TestUpdateStatus(t *testing.T) {
	m := message{
		id:        []byte("id"),
		createdAt: 9999,
		status:    StatusDeleted,
		body:      []byte("body"),
	}

	bytes := m.Encode()
	updateStatus(bytes, StatusPending)

	m.Decode(bytes)
	assert.Equal(t, uint16(StatusPending), m.status)
}
