package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseMessage(t *testing.T) {
	m := BaseMessage{
		id:        []byte("id"),
		createdAt: 9999,
		status:    StatusDeleted,
	}

	md := BaseMessage{}

	md.Decode(m.Encode())

	assert.Equal(t, m.id, md.id)
	assert.Equal(t, m.createdAt, md.createdAt)
	assert.Equal(t, uint16(StatusDeleted), md.status)
}
