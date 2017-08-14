package mq

import (
	"qiniu.com/video/builder"
)

var codecs = make(map[builder.Target]map[builder.Pattern]Codec)

// Register register codec for mq message
func Register(target builder.Target, pattern builder.Pattern, codec Codec) {
	tCodec := codecs[target]
	if tCodec == nil {
		tCodec = make(map[builder.Pattern]Codec)
		codecs[target] = tCodec
	}

	tCodec[pattern] = codec
}

// GetCodec return the mq codec by target and pattern
func GetCodec(target builder.Target, pattern builder.Pattern) Codec {
	tCodec := codecs[target]
	if tCodec == nil {
		return nil
	}
	return tCodec[pattern]
}
