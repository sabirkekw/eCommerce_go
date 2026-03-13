package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// getting logger from context
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		panic("failed to recieve logger from context")
	}
	logger.Debugw("Recieved request: ", "method", serverInfo.FullMethod, "request", req)

	// validating JWT token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no token")
	}
	token := md["authorization"]
	isValid, err := valid(token, ctx.Value("jwtSecret").(string))
	if err != nil {
		if errors.Is(err, apierrors.ErrTokenExpired) {
			return nil, status.Errorf(codes.Unauthenticated, "token expired")
		}
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	if !isValid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	logger.Debugw("Token valid", "method", serverInfo.FullMethod)

	userID, err := getUserIDFromToken(token, ctx.Value("jwtSecret").(string))
	if err != nil {
		logger.Debugw("No UserID in token", "method", serverInfo.FullMethod)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	ctx = context.WithValue(ctx, "user_id", userID)

	// handling RPC
	resp, err := handler(ctx, req)
	if err != nil {
		logger.Warnw("RPC failed: ", "method", serverInfo.FullMethod, "request", req)
		return resp, err
	}

	logger.Debugw("RPC executed: ", "method", serverInfo.FullMethod)
	return resp, nil
}

func valid(tokenSliced []string, jwtSecret string) (bool, error) {
	token := strings.Join(tokenSliced, "")
	token, _ = strings.CutPrefix(token, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, apierrors.ErrTokenExpired
		}
		return false, apierrors.ErrInvalidToken
	}

	if !parsedToken.Valid {
		return false, apierrors.ErrInvalidToken
	}
	return true, nil
}

func getUserIDFromToken(tokenSliced []string, jwtSecret string) (int64, error) {
	token := strings.Join(tokenSliced, "")
	token, _ = strings.CutPrefix(token, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}
	userID, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("couldnt find user id in token")
	}
	return int64(userID["user_id"].(float64)), nil
}
