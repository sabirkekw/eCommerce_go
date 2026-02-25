package grpcapp

import (
	"fmt"
	"net"

	grpcserver "github.com/sabirkekw/ecommerce_go/order-service/internal/grpc/order"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	Logger *zap.SugaredLogger
	Server *grpc.Server
	port   int
}

func NewGRPCServer(log *zap.SugaredLogger, port int, service grpcserver.OrderService) *GRPCApp {
	grpcServer := grpc.NewServer()

	grpcserver.Register(grpcServer, grpcserver.New(service, log))

	return &GRPCApp{
		Logger: log,
		Server: grpcServer,
		port:   port,
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
