package mq

import (
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

var codecs = make(map[target.Target]map[pattern.Pattern]Codec)

// Register register codec for mq message
func Register(t target.Target, p pattern.Pattern, codec Codec) {
	tCodec := codecs[t]
	if tCodec == nil {
		tCodec = make(map[pattern.Pattern]Codec)
		codecs[t] = tCodec
	}

	tCodec[p] = codec
}

// GetCodec return the mq codec by target and pattern
func GetCodec(target target.Target, pattern pattern.Pattern) Codec {
	tCodec := codecs[target]
	if tCodec == nil {
		return nil
	}
	return tCodec[pattern]
}
