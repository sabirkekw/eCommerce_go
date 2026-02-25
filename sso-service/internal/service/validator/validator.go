package validator

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type ValidatorService struct {
	logger    *zap.SugaredLogger
	jwtSecret string
}

func New(logger *zap.SugaredLogger, jwtSecret string) *ValidatorService {
	return &ValidatorService{
		logger:    logger,
		jwtSecret: jwtSecret,
	}
}

func (s *ValidatorService) Validate(ctx context.Context, token string) (bool, error) {
	const op = "sso.Validator.Service.Validate"

	s.logger.Infow("Validating token", "op", op)

	if token == "" {
		return false, ErrInvalidToken
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Infow("Unexpected signing method", "method", token.Header["alg"], "op", op)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			s.logger.Infow("Token expired", "op", op)
			return false, ErrTokenExpired
		}
		s.logger.Infow("Failed to parse token", "error", err, "op", op)
		return false, ErrInvalidToken
	}

	if !parsedToken.Valid {
		s.logger.Infow("Invalid token", "op", op)
		return false, ErrInvalidToken
	}
	return true, nil
}
