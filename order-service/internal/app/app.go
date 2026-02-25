package app

import (
	grpcapp "github.com/sabirkekw/ecommerce_go/order-service/internal/app/grpc"
	grpcserver "github.com/sabirkekw/ecommerce_go/order-service/internal/grpc/order"
	"github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcapp.GRPCApp
	Storage    map[string]*order.OrderData
}

func New(log *zap.SugaredLogger, port int, storage map[string]*order.OrderData, service grpcserver.OrderService) *App {
	GRPCServer := grpcapp.NewGRPCServer(log, port, service)

	return &App{
		GRPCServer: GRPCServer,
		Storage:    storage,
	}
}
