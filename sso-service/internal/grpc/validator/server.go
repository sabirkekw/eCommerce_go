package validator

import (
	"context"

	proto "github.com/sabirkekw/ecommerce_go/pkg/api/sso"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ValidatorServer struct {
	validator ValidatorService
	logger    *zap.SugaredLogger
	proto.UnimplementedValidatorServer
}

type ValidatorService interface {
	Validate(ctx context.Context, token string) (bool, error)
}

func New(validator ValidatorService, logger *zap.SugaredLogger) *ValidatorServer {
	return &ValidatorServer{
		validator: validator,
		logger:    logger,
	}
}

func Register(grpc *grpc.Server, server *ValidatorServer) {
	proto.RegisterValidatorServer(grpc, server)
}

func (s *ValidatorServer) Validate(ctx context.Context, in *proto.ValidatorRequest) (*proto.ValidatorResponse, error) {
	const op = "sso.Validator.Server.Validate"

	s.logger.Infow("Recieved Validate request", "op", op)
	isValid, err := s.validator.Validate(ctx, in.Token)
	if err != nil {
		return nil, err
	}
	return &proto.ValidatorResponse{IsValid: isValid}, nil
}
