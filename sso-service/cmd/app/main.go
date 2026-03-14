package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/sabirkekw/ecommerce_go/pkg/logger"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/app"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/config"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/repository"
	authservice "github.com/sabirkekw/ecommerce_go/sso-service/internal/service/auth"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()
	logger.Log.Infow("Logger initialized")

	cfg := config.MustLoad()
	logger.Log.Infow("Config loaded\n", "config", cfg)

	db, err := repository.ConnectToPostgres(cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Database, logger.Log)
	if err != nil {
		logger.Log.Errorw("Failed to connect to Postgres", "error", err)
		return
	}
	defer db.Close()
	logger.Log.Infow("Connected to Postgres")

	authRepo := repository.New(logger.Log, db, squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar))

	authService := authservice.New(logger.Log, authRepo, cfg.TokenTTL, cfg.JWTSecret)

	application := app.New(logger.Log, cfg.GRPC.Port, cfg.HTTP.Port, db, authService, cfg.GRPC.Timeout)
	go application.AuthGRPCServer.Run()
	logger.Log.Infow("gRPC server started", "auth_port", cfg.GRPC.Port)
	go application.AuthHTTPServer.Run()
	logger.Log.Infow("HTTP gateway started", "auth_port", cfg.HTTP.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	application.AuthGRPCServer.Stop()

	logger.Log.Infow("Server received shutdown signal, exiting")
}
