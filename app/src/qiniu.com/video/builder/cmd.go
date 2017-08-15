package builder

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

type flowRandom struct{}

func (f flowRandom) Build(video string,
	params interface{},
) ([]interface{}, error) {
	// TODO
	return nil, nil
}

func init() {
	Register(Cmd, targetFrame, patternRandom, frameRandom{})
	Register(Cmd, targetFlow, patternRandom, flowRandom{})
}
