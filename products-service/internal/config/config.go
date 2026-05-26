package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string `yaml:"env" env:"APP_ENV" env-default:"local"`
	Storage struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"postgres"`
		Port     int    `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
		Username string `yaml:"username" env:"POSTGRES_USER" env-default:"postgres"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"postgres"`
		Database string `yaml:"database" env:"POSTGRES_DB" env-default:"postgres"`
	} `yaml:"storage"`
	GRPC struct {
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"grpc"`
	HTTP struct {
		Port    int           `yaml:"port"`
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
