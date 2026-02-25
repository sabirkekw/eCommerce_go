package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sabirkekw/ecommerce_go/order-service/internal/app"
	"github.com/sabirkekw/ecommerce_go/order-service/internal/cfg"
	"github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	orderservice "github.com/sabirkekw/ecommerce_go/order-service/internal/services/order"
	"github.com/sabirkekw/ecommerce_go/pkg/logger"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()
	logger.Log.Infow("Logger initialized")

	config := cfg.MustLoad()
	logger.Log.Infow("Config loaded\n", "config", fmt.Sprintf("%+v", config))

	db := make(map[string]*order.OrderData)

	orderService := orderservice.NewService(db, logger.Log)

	application := app.New(logger.Log, config.GRPC.Port, db, orderService)

	go application.GRPCServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	application.GRPCServer.Stop()
	logger.Log.Infow("Server received shutdown signal, exiting")
}
