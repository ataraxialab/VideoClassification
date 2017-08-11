package frame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodec(t *testing.T) {
	f := Message{
		Frame: Frame{
			Index:     1,
			Label:     .123,
			ImagePath: "a/b/c",
		},
	}
	f.SetID("id")
	f.SetCreatedAt(90900)

	assert.Equal(t, "id", f.ID())
	assert.Equal(t, uint64(90900), f.CreatedAt())

	df := Message{}
	df.Decode(f.Encode())

	assert.Equal(t, f.Index, df.Index)
	assert.Equal(t, f.Label, df.Label)
	assert.Equal(t, f.ImagePath, df.ImagePath)
	assert.Equal(t, f.CreatedAt(), df.CreatedAt())
	assert.Equal(t, f.ID(), df.ID())
}
