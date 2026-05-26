package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sabirkekw/ecommerce_go/cart-service/internal/grpc/interceptor"
	cartgrpc "github.com/sabirkekw/ecommerce_go/cart-service/internal/grpc/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	Logger *zap.SugaredLogger
	Server *grpc.Server
	Port   int
}

func New(logger *zap.SugaredLogger, port int, service cartgrpc.CartService, jwtSecret string, timeout time.Duration) *GRPCApp {
	logInterceptor := func(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = context.WithValue(ctx, "logger", logger)
		ctx = context.WithValue(ctx, "jwtSecret", jwtSecret)
		ctx = context.WithValue(ctx, "timeout", timeout)
		return interceptor.LogInterceptor(ctx, req, serverInfo, handler)
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logInterceptor,
		interceptor.TimeoutInterceptor,
		interceptor.AuthInterceptor,
		interceptor.UserIDExtractorInterceptor,
	),
	)

	cartgrpc.Register(grpcServer, cartgrpc.New(service, logger))

	return &GRPCApp{
		Logger: logger,
		Server: grpcServer,
		Port:   port,
	}
}

func (s *GRPCApp) Run() {
	const op = "Cart.GRPCApp.Run"

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		s.Logger.Fatalw("failed to listen on port", "error", err, "port", s.Port, "op", op)
	}
	if err := s.Server.Serve(listener); err != nil {
		s.Logger.Fatalw("failed to start gRPC server", "error", err, "port", s.Port, "op", op)
	}
}

func (s *GRPCApp) Stop() {
	const op = "Cart.GRPCApp.Stop"
	s.Logger.Infow("Stopping gRPC server", "op", op)
	s.Server.GracefulStop()
}
