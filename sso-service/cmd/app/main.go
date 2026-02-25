package main

import (
	"github.com/sabirkekw/ecommerce_go/pkg/logger"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/config"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()
	logger.Log.Infow("Logger initialized")

	cfg := config.MustLoad()
	logger.Log.Infow("Config loaded\n", "config", cfg)

	// TODO: init postgres

	// TODO: init auth and validator servers and run them

	// TODO: graceful shutdown
}
