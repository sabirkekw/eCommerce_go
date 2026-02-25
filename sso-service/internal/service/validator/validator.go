package validator

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type ValidatorService struct {
	jwtSecret string
}

func New(jwtSecret string) *ValidatorService {
	return &ValidatorService{
		jwtSecret: jwtSecret,
	}
}

func (s *ValidatorService) Validate(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, ErrInvalidToken
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, ErrTokenExpired
		}
		return false, ErrInvalidToken
	}

	if !parsedToken.Valid {
		return false, ErrInvalidToken
	}
	return true, nil
}
