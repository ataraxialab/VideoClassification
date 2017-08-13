package builder

import "qiniu.com/video"

// Builder build the data from vido
type Builder interface {
	Builder(video string, params interface{}) ([]interface{}, error)
}

// Implement the build implementation
type Implement string

// Target build target, like cut-frame, calculate flow
type Target string

const (
	// Cmd call external command to build data
	Cmd Implement = "cmd"
	// Frame create for frame
	Frame Target = "frame"
	// Flow create for flow
	Flow Target = "flow"
)

var builders = make(map[Implement]map[Target]map[video.Pattern]Builder)

func register(impl Implement,
	target Target,
	pattern video.Pattern,
	builder Builder,
) {
	tBuilders := builders[impl]
	if tBuilders == nil {
		tBuilders = make(map[Target]map[video.Pattern]Builder)
		builders[impl] = tBuilders
	}

	pBuilders = tBuilders[pattern]
	if pBuilders == nil {
		pBuilders = make(map[video.Pattern]Builder)
		tBuilders[pattern] = pBuilders
	}
	pBuilders[target] = builder
}

// GetBuilder return the data builder
func GetBuilder(impl Implement, target Target, pattern video.Pattern) Builder {
	tBuilders := builders[impl]
	if tBuilders == nil {
		return nil
	}

	pBuilders = tBuilders[pattern]
	if pBuilders == nil {
		return nil
	}

	return pBuilders[target]
}
