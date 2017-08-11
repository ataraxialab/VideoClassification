package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseMessage(t *testing.T) {
	m := BaseMessage{
		id:        "id",
		createdAt: 9999,
	}

	md := BaseMessage{}

	md.Decode(m.Encode())

	assert.Equal(t, m.id, md.id)
	assert.Equal(t, m.createdAt, md.createdAt)
}
