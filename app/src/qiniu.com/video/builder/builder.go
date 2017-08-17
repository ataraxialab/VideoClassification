package builder

// Builder build the data from vido
type Builder interface {
	Build(params interface{}) ([]interface{}, error)
	Clean(interface{}) error
}

// Implement the build implementation
type Implement string

// Params building parameters, it is BAD design
type Params struct {
	Count  int     `json:"count"`
	Offset float32 `json:"offset"`
}

var outputRoot = "build-output"

// SetOutputRoot set build output directory
func SetOutputRoot(d string) {
	outputRoot = d
}
