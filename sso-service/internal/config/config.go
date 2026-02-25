package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string `yaml:"env" env-default:"development"`
	Storage struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
	} `yaml:"storage"`
	GRPC struct {
		Port    int
		Timeout time.Duration
	} `yaml:"grpc"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig("local.yaml", &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
