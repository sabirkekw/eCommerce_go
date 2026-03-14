package app

import (
	"database/sql"
	"time"

	authgrpcapp "github.com/sabirkekw/ecommerce_go/sso-service/internal/app/grpc"
	authhttpapp "github.com/sabirkekw/ecommerce_go/sso-service/internal/app/http"
	authgrpcserver "github.com/sabirkekw/ecommerce_go/sso-service/internal/grpc/auth"
	"go.uber.org/zap"
)

type App struct {
	AuthGRPCServer *authgrpcapp.AuthGRPCApp
	AuthHTTPServer *authhttpapp.AuthHTTPApp
	Storage        *sql.DB
}

func New(log *zap.SugaredLogger, GRPCPort int, HTTPPort int, storage *sql.DB, authService authgrpcserver.AuthService, timeout time.Duration) *App {
	authGRPCServer := authgrpcapp.NewGRPCServer(log, GRPCPort, authService, timeout)
	authHTTPServer := authhttpapp.New(log, HTTPPort, GRPCPort)

	return &App{
		AuthHTTPServer: authHTTPServer,
		AuthGRPCServer: authGRPCServer,
		Storage:        storage,
	}
}
