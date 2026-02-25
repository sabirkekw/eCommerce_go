package validatorgrpcapp

import (
	"fmt"
	"net"

	validatorgrpc "github.com/sabirkekw/ecommerce_go/sso-service/internal/grpc/validator"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ValidatorGRPCApp struct {
	Logger *zap.SugaredLogger
	Server *grpc.Server
	port   int
}

func NewGRPCServer(log *zap.SugaredLogger, port int, service validatorgrpc.ValidatorService) *ValidatorGRPCApp {
	grpcServer := grpc.NewServer()

	validatorgrpc.Register(grpcServer, validatorgrpc.New(service, log))
	return &ValidatorGRPCApp{
		Logger: log,
		Server: grpcServer,
		port:   port,
	}
}

func (s *ValidatorGRPCApp) Run() {
	const op = "Validator.grpcapp.Run"
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

func (s* ValidatorGRPCApp) Stop() {
	const op = "Validator.grpcapp.Stop"
	s.Logger.Info("Stopping gRPC server on port %d", s.port)
	s.Server.GracefulStop()
}