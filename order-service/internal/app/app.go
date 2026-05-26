package app

import (
	"database/sql"
	"time"

	grpcapp "github.com/sabirkekw/ecommerce_go/order-service/internal/app/grpc"
	httpapp "github.com/sabirkekw/ecommerce_go/order-service/internal/app/http"
	grpcserver "github.com/sabirkekw/ecommerce_go/order-service/internal/grpc/server"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcapp.GRPCApp
	HTTPServer *httpapp.HTTPApp
	Storage    *sql.DB
}

func New(log *zap.SugaredLogger, grpcport int, httpport int, storage *sql.DB, service grpcserver.OrderService, jwtSecret string, timeout time.Duration) *App {
	GRPCServer := grpcapp.NewGRPCServer(log, grpcport, service, timeout, jwtSecret)
	HTTPServer := httpapp.New(log, httpport, grpcport)

	return &App{
		GRPCServer: GRPCServer,
		HTTPServer: HTTPServer,
		Storage:    storage,
	}
}
