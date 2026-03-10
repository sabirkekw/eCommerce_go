package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env"`
	Storage struct {
		Host string `yaml:"host"`
		Port int `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"storage"`
	GRPC struct {
		Port int `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"grpc"`
	HTTP struct {
		Port int `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"http"`
	JWTSecret string `yaml:"jwt_secret"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig("./config/products-service/local.yaml", &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
