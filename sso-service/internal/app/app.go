package app

import (
	"database/sql"

	"go.uber.org/zap"
	authgrpcserver "github.com/sabirkekw/ecommerce_go/sso-service/internal/grpc/auth"
	validatorgrpcserver "github.com/sabirkekw/ecommerce_go/sso-service/internal/grpc/validator"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/app/grpc/auth"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/app/grpc/validator"
)

type App struct {
	AuthGRPCServer      *authgrpcapp.AuthGRPCApp
	ValidatorGRPCServer *validatorgrpcapp.ValidatorGRPCApp
	Storage             *sql.DB
}

func New(log *zap.SugaredLogger, authPort int, validatorPort int, storage *sql.DB, authService authgrpcserver.AuthService, validatorService validatorgrpcserver.ValidatorService) *App {

	authGRPCServer := authgrpcapp.NewGRPCServer(log, authPort, authService)
	validatorGRPCServer := validatorgrpcapp.NewGRPCServer(log, validatorPort, validatorService)

	return &App{
		AuthGRPCServer:      authGRPCServer,
		ValidatorGRPCServer: validatorGRPCServer,
		Storage:             storage,
	}
}