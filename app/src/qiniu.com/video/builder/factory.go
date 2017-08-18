package builder

import (
	"qiniu.com/video/pattern"
	"qiniu.com/video/target"
)

var builders = make(map[Implement]map[target.Target]map[pattern.Pattern]Builder)

// Register register builder
func Register(impl Implement,
	t target.Target,
	p pattern.Pattern,
	builder Builder,
) {
	tBuilders := builders[impl]
	if tBuilders == nil {
		tBuilders = make(map[target.Target]map[pattern.Pattern]Builder)
		builders[impl] = tBuilders
	}

	pBuilders := tBuilders[t]
	if pBuilders == nil {
		pBuilders = make(map[pattern.Pattern]Builder)
		tBuilders[t] = pBuilders
	}
	pBuilders[p] = builder
}

// GetBuilder return the data builder
func GetBuilder(impl Implement,
	target target.Target,
	pattern pattern.Pattern,
) Builder {
	tBuilders := builders[impl]
	if tBuilders == nil {
		return nil
	}

	pBuilders := tBuilders[target]
	if pBuilders == nil {
		return nil
	}

	return pBuilders[pattern]
}

// HasImplement checks implementation exists
func HasImplement(impl Implement) bool {
	return builders[impl] != nil
}
