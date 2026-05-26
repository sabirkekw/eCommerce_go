package cfg

import (
	"time"

	cleanenv "github.com/ilyakaznacheev/cleanenv"
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
	Kafka struct {
		Brokers string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"localhost:9092"`
		GroupID string `yaml:"group_id" env:"KAFKA_GROUP_ID" env-default:"order-service-group"`
		Topic   string `yaml:"topic" env:"KAFKA_TOPIC" env-default:"checkout-topic"`
	} `yaml:"kafka"`
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
