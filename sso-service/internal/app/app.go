package app

import (
	"database/sql"

	authgrpcapp "github.com/sabirkekw/ecommerce_go/sso-service/internal/app/grpc/auth"
	authgrpcserver "github.com/sabirkekw/ecommerce_go/sso-service/internal/grpc/auth"
	"go.uber.org/zap"
)

type App struct {
	AuthGRPCServer *authgrpcapp.AuthGRPCApp
	Storage        *sql.DB
}

func New(log *zap.SugaredLogger, Port int, storage *sql.DB, authService authgrpcserver.AuthService) *App {

	authGRPCServer := authgrpcapp.NewGRPCServer(log, Port, authService)

	return &App{
		AuthGRPCServer: authGRPCServer,
		Storage:        storage,
	}
}
