package mq

var mqs = make(map[string]MQ)

func register(name string, mq MQ) {
	mqs[name] = mq
}

// Get return the mq by the name
func Get(name string) MQ {
	q, exits := mqs[name]
	if !exits {
		return nil
	}

	return q
}
