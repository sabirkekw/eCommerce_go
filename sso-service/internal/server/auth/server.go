package auth

import (
	"context"

	proto "github.com/sabirkekw/ecommerce_go/pkg/api/sso"
)

type AuthServer struct {
	port int
	auth Auth
	proto.UnimplementedAuthServer
}

type Auth interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	Register(ctx context.Context, firstName string, lastName string, email string, password string) (userID int64, err error)
}

func New(port int, auth Auth) *AuthServer {
	return &AuthServer{
		port: port,
		auth: auth,
	}
}

func (s *AuthServer) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	uid, err := s.auth.Register(ctx, in.FirstName, in.LastName, in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterResponse{Id: uid}, nil
}

func (s *AuthServer) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.auth.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{Token: token}, nil
}
