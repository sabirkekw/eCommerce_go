package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sabirkekw/ecommerce_go/pkg/logger"
	app "github.com/sabirkekw/ecommerce_go/products-service/internal/app"
	"github.com/sabirkekw/ecommerce_go/products-service/internal/config"
	"github.com/sabirkekw/ecommerce_go/products-service/internal/database/postgres"
	"github.com/sabirkekw/ecommerce_go/products-service/internal/repository"
	"github.com/sabirkekw/ecommerce_go/products-service/internal/service"
)

func main() {

	logger.InitLogger()
	defer logger.Log.Sync()
	logger.Log.Infow("Logger initialized")

	cfg := config.MustLoad()
	logger.Log.Infow("config initialized", "cfg: ", cfg)

	db, err := postgres.ConnectToPostgres(cfg)
	if err != nil {
		logger.Log.Fatalw("failed to connect to postgres", "error: ", err)
	}

	productsRepository := repository.New(db, logger.Log)

	productsService := service.New(productsRepository, logger.Log)

	application := app.New(logger.Log, cfg.HTTP.Port, cfg.GRPC.Port, productsService, cfg.JWTSecret, cfg.GRPC.Timeout)

	go application.GRPCApp.Run()
	logger.Log.Infow("Products gRPC server started", "port", cfg.GRPC.Port)
	go application.HTTPApp.Run()
	logger.Log.Infow("Products HTTP gateway started", "port", cfg.HTTP.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	application.GRPCApp.Stop()
}
