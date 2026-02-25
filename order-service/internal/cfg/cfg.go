package cfg

import (
	"time"

	cleanenv "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`
	// Storage struct {
	// 	Host     string `yaml:"host"`
	// 	Port     int    `yaml:"port"`
	// 	Username string `yaml:"username"`
	// 	Password string `yaml:"password"`
	// 	Database string `yaml:"database"`
	// } `yaml:"storage"`
	GRPC struct {
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"grpc"`
}

// done: implement config loading and validation
// done: define config struct with necessary fields (e.g. server port, database connection string, etc.)
// done: support loading config from environment variables and/or config files (e.g. YAML, JSON, etc.)

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig("./config/order-service/local.yaml", &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
