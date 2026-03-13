package interceptor

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

func AuthInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// getting logger from context
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		panic("failed to recieve logger from context")
	}
	logger.Debugw("Recieved request: ", "method", serverInfo.FullMethod, "request", req)

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Warnw("RPC failed: ", "method", serverInfo.FullMethod, "request", req)
		return resp, err
	}

	logger.Debugw("RPC executed: ", "method", serverInfo.FullMethod)
	return resp, nil
}
