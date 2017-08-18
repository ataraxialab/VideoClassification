package config

// Builder building related configurations
type Builder struct {
	MaxRetainMessageCount int    `json:"max_retain_message_count"`
	CheckPeriod           int    `json:"check_period"`
	OutputRoot            string `json:"output_root"`
	VideoRoot             string `json:"video_root"`
	TrainLabelFile        string `json:"train_label"`
	ValLabelFile          string `json:"val_label"`
}

// MQ mq related configurations
type MQ struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// Config root structure
type Config struct {
	Builder Builder `json:"builder"`
	MQ      MQ      `json:"mq"`
}
