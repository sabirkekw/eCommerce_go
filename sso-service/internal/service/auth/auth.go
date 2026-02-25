package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/models"
	"github.com/sabirkekw/ecommerce_go/sso-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, firstName string, lastName string, email string, password_hash []byte) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	userRepo  UserRepo
	tokenTTL  time.Duration
	jwtSecret string
}

func New(userRepo UserRepo, tokenTTL time.Duration, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenTTL:  tokenTTL,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, firstName string, lastName string, email string, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	id, err := s.userRepo.CreateUser(ctx, firstName, lastName, email, hash)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNoUser) {
			return "", err
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(existingUser.PassHash, []byte(password))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": existingUser.ID,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(os.Getenv("JWT_SECRET"))
}
