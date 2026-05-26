package grpcapp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sabirkekw/ecommerce_go/order-service/internal/grpc/interceptor"
	grpcserver "github.com/sabirkekw/ecommerce_go/order-service/internal/grpc/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	Logger    *zap.SugaredLogger
	Server    *grpc.Server
	port      int
	timeout   time.Duration
	JWTSecret string
}

func NewGRPCServer(log *zap.SugaredLogger, port int, service grpcserver.OrderService, timeout time.Duration, jwtSecret string) *GRPCApp {
	logInterceptor := func(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = context.WithValue(ctx, "logger", log)
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

	grpcserver.Register(grpcServer, grpcserver.New(service, log))

	return &GRPCApp{
		Logger:    log,
		Server:    grpcServer,
		port:      port,
		JWTSecret: jwtSecret,
	}
}

func (a *GRPCApp) Run() {
	const op = "grpcapp.Run"

	a.Logger.Infow("Starting gRPC server", "port", a.port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.Logger.Errorw("Failed to listen on port", "error", err, "op", op)
		return
	}

	if err := a.Server.Serve(listener); err != nil {
		a.Logger.Errorw("Failed to start gRPC server", "error", err, "op", op)
		return
	}
}

func (a *GRPCApp) Stop() {
	const op = "grpcapp.Stop"
	a.Logger.Infow("Stopping gRPC server", "op", op)
	a.Server.GracefulStop()
}
