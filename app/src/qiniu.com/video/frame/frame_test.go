package frame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodec(t *testing.T) {
	f := Frame{
		Index:     1,
		Label:     .123,
		ImagePath: "a/b/c",
		Timestamp: 1990,
		UID:       "id",
	}

	df := Frame{}
	df.Decode(f.Encode())

	assert.Equal(t, f.Index, df.Index)
	assert.Equal(t, f.Label, df.Label)
	assert.Equal(t, f.ImagePath, df.ImagePath)
	assert.Equal(t, f.Timestamp, df.Timestamp)
	assert.Equal(t, f.UID, df.UID)
}
