package builder

// Builder build the data from vido
type Builder interface {
	Build(video string, params interface{}) ([]interface{}, error)
	Clean(interface{}) error
}

// Implement the build implementation
type Implement string

// Target build target, like cut-frame, calculate flow
type Target string

const (
	// TargetFrame create for frame
	targetFrame Target = "frame"
	// TargetFlow create for flow
	targetFlow Target = "flow"
	// TargetUnknown unknow target
	targetUnknown Target = "unknown"
)

// GetTarget return the target by string
func GetTarget(t string) Target {
	switch t {
	case string(targetFrame):
		return targetFrame
	case string(targetFlow):
		return targetFlow
	default:
		return targetUnknown
	}
}

// IsValid checks target
func (t *Target) IsValid() bool {
	return *t == targetFrame || *t == targetFlow
}

// Pattern generate patterns
type Pattern string

const (
	// PatternRandom cut one frame from the video randomly
	patternRandom Pattern = "random"
	// PatternUnknown unknow target
	patternUnknown Pattern = "unknown"
)

// GetPattern return pattern by string
func GetPattern(p string) Pattern {
	switch p {
	case string(patternRandom):
		return patternRandom
	default:
		return patternUnknown
	}
}

// IsValid checks pattern
func (p *Pattern) IsValid() bool {
	return *p == patternRandom
}

// Params building parameters, it is BAD design
type Params struct {
	Count  int     `json:"count"`
	Offset float32 `json:"offset"`
}
