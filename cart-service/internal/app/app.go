package app

import (
	"time"

	grpcapp "github.com/sabirkekw/ecommerce_go/cart-service/internal/app/grpc"
	httpapp "github.com/sabirkekw/ecommerce_go/cart-service/internal/app/http"
	cartservice "github.com/sabirkekw/ecommerce_go/cart-service/internal/grpc/server"
	"go.uber.org/zap"
)

type App struct {
	GRPCApp grpcapp.GRPCApp
	HTTPApp httpapp.HTTPApp
}

func New(logger *zap.SugaredLogger, grpcPort int, httpPort int, cartService cartservice.CartService, jwtSecret string, timeout time.Duration) *App {
	grpcApp := grpcapp.New(logger, grpcPort, cartService, jwtSecret, timeout)
	httpApp := httpapp.New(logger, httpPort, grpcPort)

	return &App{
		GRPCApp: *grpcApp,
		HTTPApp: *httpApp,
	}
}
