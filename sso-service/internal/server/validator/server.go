package validator

import (
	"context"

	proto "github.com/sabirkekw/ecommerce_go/pkg/api/sso"
)

type ValidatorServer struct {
	port      int
	validator Validator
	proto.UnimplementedValidatorServer
}

type Validator interface {
	Validate(ctx context.Context, token string) (bool, error)
}

func New(port int, validator Validator) *ValidatorServer {
	return &ValidatorServer{
		port:      port,
		validator: validator,
	}
}

func (s *ValidatorServer) Validate(ctx context.Context, in *proto.ValidatorRequest) (*proto.ValidatorResponse, error) {
	isValid, err := s.validator.Validate(ctx, in.Token)
	if err != nil {
		return nil, err
	}
	return &proto.ValidatorResponse{IsValid: isValid}, nil
}
