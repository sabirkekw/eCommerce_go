package http

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gw "github.com/sabirkekw/ecommerce_go/pkg/api/products"
)

type HTTPApp struct {
	Logger *zap.SugaredLogger
	HTTPPort int
	GRPCPort int
	Router *runtime.ServeMux
}

func New(logger *zap.SugaredLogger, httpport int, grpcport int) *HTTPApp {
	router := runtime.NewServeMux()
	return &HTTPApp{
		Logger: logger,
		HTTPPort: httpport,
		GRPCPort: grpcport,
		Router: router,
	}
}

func (s *HTTPApp) Run() {
	const op = "Products.HTTPApp.Run"
	ctx := context.Background()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterProductsServiceHandlerFromEndpoint(ctx, s.Router, fmt.Sprintf("localhost:%d", s.GRPCPort), opts)
	if err != nil {
		s.Logger.Errorw("failed to register gRPC gateway", "error", err.Error(), "op", op)
		panic(err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", s.HTTPPort), s.Router)
	if err != nil {
		s.Logger.Errorw("failed to run gRPC gateway", "error", err.Error(), "op", op)
		panic(err)
	}
}
