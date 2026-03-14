package app

import (
	"database/sql"
	"time"

	grpcapp "github.com/sabirkekw/ecommerce_go/order-service/internal/app/grpc"
	grpcserver "github.com/sabirkekw/ecommerce_go/order-service/internal/grpc/server"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcapp.GRPCApp
	Storage    *sql.DB
}

func New(log *zap.SugaredLogger, port int, storage *sql.DB, service grpcserver.OrderService, jwtSecret string, timeout time.Duration) *App {
	GRPCServer := grpcapp.NewGRPCServer(log, port, service, timeout, jwtSecret)

	return &App{
		GRPCServer: GRPCServer,
		Storage:    storage,
	}
}
