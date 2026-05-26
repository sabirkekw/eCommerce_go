package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sabirkekw/ecommerce_go/cart-service/internal/app"
	"github.com/sabirkekw/ecommerce_go/cart-service/internal/config"
	"github.com/sabirkekw/ecommerce_go/cart-service/internal/database/postgres"
	"github.com/sabirkekw/ecommerce_go/cart-service/internal/database/redis"
	productsclient "github.com/sabirkekw/ecommerce_go/cart-service/internal/grpc/client"
	postgresrepo "github.com/sabirkekw/ecommerce_go/cart-service/internal/repository/postgres"
	redisrepo "github.com/sabirkekw/ecommerce_go/cart-service/internal/repository/redis"
	service "github.com/sabirkekw/ecommerce_go/cart-service/internal/service/cart"
	"github.com/sabirkekw/ecommerce_go/cart-service/internal/service/messaging"
	"github.com/sabirkekw/ecommerce_go/pkg/logger"
)

func main() {
	logger.InitLogger()
	logger.Log.Infow("Starting cart service")
	cfg := config.MustLoad()
	logger.Log.Infow("Config initialized", "cfg", fmt.Sprintf("%+v", cfg))

	postgres_db := postgres.ConnectToPostgres(cfg)
	logger.Log.Infow("connected to PostgreSQL")
	redis_db := redis.ConnectToRedis(cfg)
	logger.Log.Infow("connected to Redis")

	postgresRepo := postgresrepo.New(postgres_db, logger.Log)
	redisRepo := redisrepo.New(redis_db, logger.Log)

	productsClient := productsclient.New(logger.Log, 50052)

	kafkaProducer := messaging.New(logger.Log, "checkout-topic")
	defer kafkaProducer.Close()

	service := service.New(postgresRepo, redisRepo, productsClient, kafkaProducer, logger.Log)

	application := app.New(logger.Log, cfg.GRPC.Port, cfg.HTTP.Port, service, cfg.JWTSecret, cfg.GRPC.Timeout)
	go application.GRPCApp.Run()
	logger.Log.Infow("Starting gRPC server")
	go application.HTTPApp.Run()
	logger.Log.Infow("Starting HTTP gateway")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	application.GRPCApp.Stop()
	logger.Log.Infow("gracefully stopped gRPC server")
}
