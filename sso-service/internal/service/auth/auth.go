package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/models"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, firstName string, lastName string, email string, password_hash []byte) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	logger    *zap.SugaredLogger
	userRepo  UserRepo
	tokenTTL  time.Duration
	jwtSecret string
}

func New(logger *zap.SugaredLogger, userRepo UserRepo, tokenTTL time.Duration, jwtSecret string) *AuthService {
	return &AuthService{
		logger:    logger,
		userRepo:  userRepo,
		tokenTTL:  tokenTTL,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, firstName string, lastName string, email string, password string) (int64, error) {
	const op = "sso.Auth.Service.Register"

	s.logger.Infow("Registering new user", "op", op)

	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		s.logger.Infow("User already exists", "email", email, "op", op)
		return 0, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Infow("Failed to make password hash", "password", password, "op", op)
		return 0, errors.New("failed to hash password")
	}
	id, err := s.userRepo.CreateUser(ctx, firstName, lastName, email, hash)
	if err != nil {
		s.logger.Infow("Failed to add user to database", "op", op)
		return 0, err
	}

	s.logger.Infow("Successfuly registered user", "op", op)
	return id, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "sso.Auth.Service.Login"

	s.logger.Infow("Logging in user", "email", email, "op", op)
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNoUser) {
			s.logger.Infow("User not found", "email", email, "op", op)
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(existingUser.PassHash, []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": existingUser.ID,
		"exp":     time.Now().Add(s.tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s.logger.Infow("Successfuly logged in user", "email", email, "op", op)
	return token.SignedString([]byte(s.jwtSecret))
}
