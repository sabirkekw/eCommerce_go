package auth

import (
	"context"

	proto "github.com/sabirkekw/ecommerce_go/pkg/api/sso"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	auth   AuthService
	logger *zap.SugaredLogger
	proto.UnimplementedAuthServer
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	Register(ctx context.Context, firstName string, lastName string, email string, password string) (userID int64, err error)
}

func New(auth AuthService, logger *zap.SugaredLogger) *AuthServer {
	return &AuthServer{
		auth:   auth,
		logger: logger,
	}
}

func Register(grpc *grpc.Server, server *AuthServer) {
	proto.RegisterAuthServer(grpc, server)
}

func (s *AuthServer) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	const op = "sso.Auth.Server.Register"

	s.logger.Infow("Recieved Register request", "data", in, "op", op)

	if in.FirstName == "" || in.LastName == "" || in.Email == "" || in.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "all fields are required")
	}

	uid, err := s.auth.Register(ctx, in.FirstName, in.LastName, in.Email, in.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", in.Email)
		}
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}
	return &proto.RegisterResponse{Id: uid}, nil
}

func (s *AuthServer) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	const op = "sso.Auth.Server.Login"

	if in.Email == "" || in.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	s.logger.Infow("Recieved Login request", "op", op)
	token, err := s.auth.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials")
	}
	return &proto.LoginResponse{Token: token}, nil
}
