package flow

import (
	"encoding/binary"
	"math"

	"qiniu.com/video/mq"
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

var endian = binary.BigEndian

// Flow store the flow information of video
type Flow struct {
	Index     uint32
	Label     float32
	ImagePath string
}

type flowCodec struct{}

// Encode the flow to byte slice
func (c flowCodec) Encode(flow interface{}) []byte {
	f := flow.(Flow)
	imagePathLen := len(f.ImagePath)
	bytes := make([]byte, 4+4+2+imagePathLen)
	endian.PutUint32(bytes, f.Index)
	endian.PutUint32(bytes[4:], math.Float32bits(f.Label))
	endian.PutUint16(bytes[4+4:], uint16(imagePathLen))
	copy(bytes[4+4+2:], f.ImagePath)

	return bytes
}

// Decode fill the object from the raw bytes
func (c flowCodec) Decode(bytes []byte) interface{} {
	f := Flow{}
	f.Index = endian.Uint32(bytes)
	f.Label = math.Float32frombits(endian.Uint32(bytes[4:]))
	imagePathLen := endian.Uint16(bytes[4+4:])
	f.ImagePath = string(bytes[4+4+2 : 4+4+2+imagePathLen])

	return f
}

func init() {
	mq.Register(target.Flow, pattern.Random, flowCodec{})
}

var fc flowCodec

// Decoder return flow decoder
func Decoder() mq.Decoder {
	return fc
}

// Encoder return flow encoder
func Encoder() mq.Encoder {
	return fc
}
