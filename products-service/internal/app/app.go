package app

import (
	grpcapp "github.com/sabirkekw/ecommerce_go/products-service/internal/app/grpc"
	httpapp "github.com/sabirkekw/ecommerce_go/products-service/internal/app/http"
	productsservice "github.com/sabirkekw/ecommerce_go/products-service/internal/grpc/server"
	"go.uber.org/zap"
)

type App struct {
	GRPCApp grpcapp.GRPCApp
	HTTPApp httpapp.HTTPApp
}

func New(logger *zap.SugaredLogger, HTTPPort int, GRPCPort int, productsService productsservice.ProductsService, jwtSecret string) *App {
	productsGRPCServer := grpcapp.New(logger, GRPCPort, productsService, jwtSecret)
	productsHTTPGateway := httpapp.New(logger, HTTPPort, GRPCPort)

	return &App{
		GRPCApp: *productsGRPCServer,
		HTTPApp: *productsHTTPGateway,
	}
}
