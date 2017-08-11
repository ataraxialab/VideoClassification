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
	Timestamp uint64
	UID       string
}

// ID return the id of frame
func (f *Frame) ID() string {
	return f.UID
}

// Unixtime the frame generated time
func (f *Frame) Unixtime() uint64 {
	return f.Timestamp
}

// Encode the frame to byte slice
func (f *Frame) Encode() []byte {
	imagePathLen, idLen := len(f.ImagePath), len(f.UID)
	bytes := make([]byte, 4+4+2+imagePathLen+8+2+idLen)

	endian.PutUint32(bytes, f.Index)
	endian.PutUint32(bytes[4:], math.Float32bits(f.Label))
	endian.PutUint16(bytes[4+4:], uint16(imagePathLen))
	copy(bytes[4+4+2:], f.ImagePath)
	endian.PutUint64(bytes[4+4+2+imagePathLen:], f.Timestamp)
	copy(bytes[4+4+2+imagePathLen+8+2:], f.UID)

	return bytes
}

// Decode fill the object from the raw bytes
func (f *Frame) Decode(bytes []byte) {
	f.Index = endian.Uint32(bytes)
	f.Label = math.Float32frombits(endian.Uint32(bytes[4:]))
	imagePathLen := endian.Uint16(bytes[4+4:])
	f.ImagePath = string(bytes[4+4+2 : 4+4+2+imagePathLen])
	f.Timestamp = endian.Uint64(bytes[4+4+2+imagePathLen:])
	f.UID = string(bytes[4+4+2+imagePathLen+8+2:])
}
