package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageEx(t *testing.T) {
	m := MessageEx{
		ID:        []byte("id"),
		CreatedAt: 9999,
		Status:    StatusDeleted,
		Body:      []byte("Body"),
	}

	md := MessageEx{}

	md.Decode(m.Encode())

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

	bytes := m.Encode()
	updateStatus(bytes, StatusPending)

	m.Decode(bytes)
	assert.Equal(t, uint16(StatusPending), m.Status)
}
