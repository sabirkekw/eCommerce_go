package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LogInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	logger := mustLogger(ctx)
	logger.Infow("Recieved request: ", "method", serverInfo.FullMethod, "request", req)

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Infow("RPC failed: ", "method", serverInfo.FullMethod, "request", req)
		return resp, err
	}

	logger.Infow("RPC executed: ", "method", serverInfo.FullMethod)
	return resp, nil
}

func UserIDExtractorInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	logger := mustLogger(ctx)

	token, err := tokenFromContext(ctx)
	if err != nil {
		logger.Debugw("failed to recieve token from context")
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	userID, err := getUserIDFromToken(token, mustJWTSecret(ctx))
	if err != nil {
		logger.Debugw("No UserID in token", "method", serverInfo.FullMethod)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	ctx = context.WithValue(ctx, "user_id", userID)

	resp, err := handler(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func AuthInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	logger := mustLogger(ctx)

	token, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "no token")
	}

	isValid, err := validateToken(token, mustJWTSecret(ctx))
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

	// keep token around for downstream interceptors
	ctx = context.WithValue(ctx, "token", token)

	resp, err := handler(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func TimeoutInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	timeout, ok := ctx.Value("timeout").(time.Duration)
	if !ok {
		timeout = 5 * time.Second // default timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct {
		resp any
		err  error
	}, 1)
	go func() {
		resp, err := handler(ctx, req)
		done <- struct {
			resp any
			err  error
		}{resp, err}
	}()

	select {
	case <-ctx.Done():
		return nil, status.Errorf(codes.DeadlineExceeded, "timeout limit exceeded")
	case response := <-done:
		return response.resp, response.err
	}
}

func mustLogger(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		panic("failed to recieve logger from context")
	}
	return logger
}

func mustJWTSecret(ctx context.Context) string {
	secret, ok := ctx.Value("jwtSecret").(string)
	if !ok {
		panic("failed to recieve jwt secret from context")
	}
	return secret
}

func tokenFromContext(ctx context.Context) ([]string, error) {
	if token, ok := ctx.Value("token").([]string); ok && len(token) > 0 {
		return token, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	token := md["authorization"]
	if len(token) == 0 {
		return nil, errors.New("missing authorization header")
	}

	return token, nil
}

func validateToken(tokenSliced []string, jwtSecret string) (bool, error) {
	token := normalizeToken(tokenSliced)
	parsedToken, err := parseToken(token, jwtSecret)
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
	token := normalizeToken(tokenSliced)
	parsedToken, err := parseToken(token, jwtSecret)
	if err != nil {
		return 0, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("couldnt find user id in token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("couldnt find user id in token")
	}

	return int64(userID), nil
}

func normalizeToken(tokenSliced []string) string {
	token := strings.Join(tokenSliced, "")
	return strings.TrimPrefix(token, "Bearer ")
}

func parseToken(token, jwtSecret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
}
