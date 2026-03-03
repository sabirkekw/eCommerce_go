package client

import (
	"context"
	"fmt"
	"strings"

	protosso "github.com/sabirkekw/ecommerce_go/pkg/api/sso"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type GRPCClient struct {
	Client protosso.ValidatorClient
	port   int
	logger *zap.SugaredLogger
}

func NewGRPCClient(port int, logger *zap.SugaredLogger) *GRPCClient {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorw("Failed to create gRPC client", "error", err)
		return nil
	}

	c := protosso.NewValidatorClient(conn)

	return &GRPCClient{
		Client: c,
		port:   port,
		logger: logger,
	}
}

func (s *GRPCClient) SendTokenToValidate(ctx context.Context, token []string) (bool, error) {
	const op = "Order.Client.ValidateToken"
	s.logger.Infow("Sending token to SSO for validation", "op", op)

	joinedToken := strings.Join(token, "")

	resp, err := s.Client.Validate(ctx, &protosso.ValidatorRequest{Token: joinedToken})
	if err != nil {
		s.logger.Infow("failed to validate token", "op", op)
		return false, status.Errorf(codes.Internal, "failed to validate token")
	}

	return resp.IsValid, nil
}
