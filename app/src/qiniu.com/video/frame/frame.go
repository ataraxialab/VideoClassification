package frame

import (
	"encoding/binary"
	"math"

	"qiniu.com/video/mq"
)

var endian = binary.BigEndian

// Frame store the frame information of video
type Frame struct {
	Index     uint32
	Label     float32
	ImagePath string
}

// Message frame message
type Message struct {
	mq.BaseMessage
	Frame
}

// Encode the frame to byte slice
func (f *Message) Encode() []byte {
	baseBytes := f.BaseMessage.Encode()

	imagePathLen, baseLen := len(f.ImagePath), len(baseBytes)
	bytes := make([]byte, 4+4+2+imagePathLen+baseLen)
	endian.PutUint32(bytes, f.Index)
	endian.PutUint32(bytes[4:], math.Float32bits(f.Label))
	endian.PutUint16(bytes[4+4:], uint16(imagePathLen))
	copy(bytes[4+4+2:], f.ImagePath)
	copy(bytes[4+4+2+imagePathLen:], baseBytes)

	return bytes
}

// Decode fill the object from the raw bytes
func (f *Message) Decode(bytes []byte) {
	f.Index = endian.Uint32(bytes)
	f.Label = math.Float32frombits(endian.Uint32(bytes[4:]))
	imagePathLen := endian.Uint16(bytes[4+4:])
	f.ImagePath = string(bytes[4+4+2 : 4+4+2+imagePathLen])

	f.BaseMessage.Decode(bytes[4+4+2+imagePathLen:])
}
