package builder

const (
	// Cmd call external command to build data
	Cmd Implement = "cmd"
)

type frameRandom struct{}

func (f frameRandom) Builder(params interface{}) ([]interface{}, error) {
	// TODO
	return nil
}

type flowRandom struct{}

func (f flowRandom) Builder(params interface{}) ([]interface{}, error) {
	// TODO
	return nil
}

func init() {
	register(Cmd, Frame, PatternRandom, frameRandom{})
	register(Cmd, Flow, PatternRandom, flowRandom{})
}
