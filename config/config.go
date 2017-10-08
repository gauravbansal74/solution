package config

type Config struct {
	// Solution App config
	ListenPort int32  `envconfig:"listen_port"`
	ListenHost string `envconfig:"listen_host"`
	Debug      bool   `envconfig:"debug"`
	SystemLog  string `envconfig:"system_log_code"`
	SystemCode string `envconfig:"system_code"`
	SystemInfo string `envconfig:"system_info"`

	// Config for Redis
	RedisHost  string `envconfig:"redis_host"`
	RedisPort  int    `envconfig:"redis_port"`
	RedisDB    int    `envconfig:"redis_db"`
	RedisQueue string `envconfig:"redis_queue"`
	// Google API key
	GoogleApiKey string `envconfig:"google_api_key"`
}
