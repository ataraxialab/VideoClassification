package target

// Target build target, like cut-frame, calculate flow
type Target string

const (
	// Frame frame
	Frame Target = "frame"
	// Flow flow
	Flow Target = "flow"
	// Unknown unknow target
	Unknown Target = "unknown"
)

// GetTarget return the target by string
func GetTarget(t string) Target {
	switch t {
	case string(Frame):
		return Frame
	case string(Flow):
		return Flow
	default:
		return Unknown
	}
}

// IsValid checks target
func (t *Target) IsValid() bool {
	return *t == Frame || *t == Flow
}
