package builder

var builders = make(map[Implement]map[Target]map[Pattern]Builder)

// Register register builder
func Register(impl Implement,
	target Target,
	pattern Pattern,
	builder Builder,
) {
	tBuilders := builders[impl]
	if tBuilders == nil {
		tBuilders = make(map[Target]map[Pattern]Builder)
		builders[impl] = tBuilders
	}

	pBuilders := tBuilders[target]
	if pBuilders == nil {
		pBuilders = make(map[Pattern]Builder)
		tBuilders[target] = pBuilders
	}
	pBuilders[pattern] = builder
}

// GetBuilder return the data builder
func GetBuilder(impl Implement,
	target Target,
	pattern Pattern,
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
