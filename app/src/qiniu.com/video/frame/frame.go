package frame

import (
	"encoding/binary"
	"math"

	"qiniu.com/video/mq"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

var endian = binary.BigEndian

// Frame store the frame information of video
type Frame struct {
	Index     uint32
	Label     float32
	ImagePath string
}

type frameCodec struct{}

// Encode the frame to byte slice
func (c frameCodec) Encode(frame interface{}) []byte {
	f := frame.(Frame)
	imagePathLen := len(f.ImagePath)
	bytes := make([]byte, 4+4+2+imagePathLen)
	endian.PutUint32(bytes, f.Index)
	endian.PutUint32(bytes[4:], math.Float32bits(f.Label))
	endian.PutUint16(bytes[4+4:], uint16(imagePathLen))
	copy(bytes[4+4+2:], f.ImagePath)

	return bytes
}

// Decode fill the object from the raw bytes
func (c frameCodec) Decode(bytes []byte) interface{} {
	f := Frame{}
	f.Index = endian.Uint32(bytes)
	f.Label = math.Float32frombits(endian.Uint32(bytes[4:]))
	imagePathLen := endian.Uint16(bytes[4+4:])
	f.ImagePath = string(bytes[4+4+2 : 4+4+2+imagePathLen])

	return f
}

func init() {
	mq.Register(target.Frame, pattern.Random, frameCodec{})
}

var fc frameCodec

// Decoder return frame decoder
func Decoder() mq.Decoder {
	return fc
}

// Encoder return frame encoder
func Encoder() mq.Encoder {
	return fc
}
