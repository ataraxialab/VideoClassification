package frame

import (
	"encoding/binary"
	"math"
)

var endian = binary.BigEndian

// Frame store the frame information of video
type Frame struct {
	Index     uint32
	Label     float32
	ImagePath string
}

// Encode the frame to byte slice
func (f *Frame) Encode() []byte {
	imagePathLen := len(f.ImagePath)
	bytes := make([]byte, 4+4+2+imagePathLen)
	endian.PutUint32(bytes, f.Index)
	endian.PutUint32(bytes[4:], math.Float32bits(f.Label))
	endian.PutUint16(bytes[4+4:], uint16(imagePathLen))
	copy(bytes[4+4+2:], f.ImagePath)

	return bytes
}

// Decode fill the object from the raw bytes
func (f *Frame) Decode(bytes []byte) {
	f.Index = endian.Uint32(bytes)
	f.Label = math.Float32frombits(endian.Uint32(bytes[4:]))
	imagePathLen := endian.Uint16(bytes[4+4:])
	f.ImagePath = string(bytes[4+4+2 : 4+4+2+imagePathLen])
}
