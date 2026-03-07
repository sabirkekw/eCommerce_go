package http

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gw "github.com/sabirkekw/ecommerce_go/pkg/api/sso"
)

type AuthHTTPApp struct {
	Logger   *zap.SugaredLogger
	Router   *runtime.ServeMux
	Port     int
	GRPCPort int
}

func New(logger *zap.SugaredLogger, port int, grpcPort int) *AuthHTTPApp {
	router := runtime.NewServeMux()
	return &AuthHTTPApp{
		Logger:   logger,
		Router:   router,
		Port:     port,
		GRPCPort: grpcPort,
	}
}

func (s *AuthHTTPApp) Run() {
	const op = "Auth.HTTPApp.Run"
	ctx := context.Background()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterAuthHandlerFromEndpoint(ctx, s.Router, fmt.Sprintf("localhost:%d", s.GRPCPort), opts)
	if err != nil {
		s.Logger.Errorw("failed to register gRPC gateway", "error", err.Error(), "op", op)
		panic(err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Router)
	if err != nil {
		s.Logger.Errorw("failed to run gRPC gateway", "error", err.Error(), "op", op)
		panic(err)
	}
}
