package frame

import (
	"testing"

	"qiniu.com/video/mq"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"

	"github.com/stretchr/testify/assert"
)

func TestCodec(t *testing.T) {
	f := Frame{
		Index:     1,
		Label:     .123,
		ImagePath: "a/b/c",
	}

	df := Decoder().Decode(Encoder().Encode(f)).(Frame)

	assert.Equal(t, f.Index, df.Index)
	assert.Equal(t, f.Label, df.Label)
	assert.Equal(t, f.ImagePath, df.ImagePath)

	assert.NotNil(t, mq.GetCodec(target.Frame, pattern.Random))
}
