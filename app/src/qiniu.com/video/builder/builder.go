package builder

// Builder build the data from vido
type Builder interface {
	Build(video string, params interface{}) ([]interface{}, error)
}

// Pattern generate patterns
type Pattern string

// Implement the build implementation
type Implement string

// Target build target, like cut-frame, calculate flow
type Target string

const (
	// Frame create for frame
	Frame Target = "frame"
	// Flow create for flow
	Flow Target = "flow"
)

const (
	// PatternRandom cut one frame from the video randomly
	PatternRandom Pattern = "random"
)
