package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string `yaml:"env" env-default:"development"`
	Storage struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"storage"`
	GRPC struct {
		AuthPort      int           `yaml:"auth_port"`
		ValidatorPort int           `yaml:"validator_port"`
		Timeout       time.Duration `yaml:"timeout"`
	} `yaml:"grpc"`
	JWTSecret string        `yaml:"jwt_secret"`
	TokenTTL  time.Duration `yaml:"token_ttl"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig("./config/sso-service/local.yaml", &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
