package config

// Builder building related configurations
type Builder struct {
	MaxRetainMessageCount int `json:"max_retain_message_count"`
	CheckPeriod           int `json:"check_period"`
}

// MQ mq related configurations
type MQ struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// HTTP http related configurations
type HTTP struct {
	Port int `json:"port"`
}

// Config root structure
type Config struct {
	Builder Builder `json:"builder"`
	MQ      MQ      `json:"mq"`
	HTTP    HTTP    `json:"http"`
}
