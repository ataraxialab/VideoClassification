package pattern

// Pattern generate patterns
type Pattern string

const (
	// Random cut one frame from the video randomly
	Random Pattern = "random"
	// Unknown unknow target
	Unknown Pattern = "unknown"
)

// GetPattern return pattern by string
func GetPattern(p string) Pattern {
	switch p {
	case string(Random):
		return Random
	default:
		return Unknown
	}
}

// IsValid checks pattern
func (p *Pattern) IsValid() bool {
	return *p == Random
}
