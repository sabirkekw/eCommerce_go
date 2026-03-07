package authgrpcapp

import (
	"fmt"
	"net"

	authgrpc "github.com/sabirkekw/ecommerce_go/sso-service/internal/grpc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AuthGRPCApp struct {
	Logger *zap.SugaredLogger
	Server *grpc.Server
	port   int
}

func NewGRPCServer(log *zap.SugaredLogger, port int, service authgrpc.AuthService) *AuthGRPCApp {
	grpcServer := grpc.NewServer()

	authgrpc.Register(grpcServer, authgrpc.New(service, log))
	return &AuthGRPCApp{
		Logger: log,
		Server: grpcServer,
		port:   port,
	}
}

func (s *AuthGRPCApp) Run() {
	const op = "Auth.grpcapp.Run"
	s.Logger.Info("Starting gRPC server on port %d", s.port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.Logger.Error("Failed to start gRPC server: %v", err)
		return
	}
	err = s.Server.Serve(listener)
	if err != nil {
		s.Logger.Error("Failed to serve gRPC server: %v", err)
	}
}

func (s *AuthGRPCApp) Stop() {
	const op = "Auth.grpcapp.Stop"
	s.Logger.Info("Stopping gRPC server on port %d", s.port)
	s.Server.GracefulStop()
}