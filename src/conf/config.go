package conf

// 配置文件
type Config struct {
	KafkaConfig	`ini:"kafka"`
	CollectConfig  `ini:"collect"`
}

// kafka配置
type KafkaConfig struct {
	Address string	`ini:"address"`
	Topic string `ini:"topic"`
	ChanSize int `ini:"chan_size"`
}

// 日志文件路径配置
type CollectConfig struct {
	LogfilePath string `ini:"logfile_path"`
}
