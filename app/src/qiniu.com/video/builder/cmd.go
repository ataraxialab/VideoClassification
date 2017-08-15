package builder

import (
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

const (
	// Cmd call external command to build data
	Cmd Implement = "cmd"
)

type frameRandom struct{}

func (f frameRandom) Build(video string,
	params interface{},
) ([]interface{}, error) {
	// TODO
	return nil, nil
}
func (f frameRandom) Clean(result interface{}) error {
	return nil
}

type flowRandom struct{}

func (f flowRandom) Build(video string,
	params interface{},
) ([]interface{}, error) {
	// TODO
	return nil, nil
}
func (f flowRandom) Clean(result interface{}) error {
	return nil
}

func init() {
	Register(Cmd, target.Frame, pattern.Random, frameRandom{})
	Register(Cmd, target.Flow, pattern.Random, flowRandom{})
}
