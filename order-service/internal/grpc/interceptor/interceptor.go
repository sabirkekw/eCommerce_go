package interceptor

import (
	"context"

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
	logger.Infow("Recieved request: ", "method", serverInfo.FullMethod, "request", req)

	// validating JWT token
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no token")
	}
	token := md["authorization"]
	if !valid(token) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// handling RPC
	resp, err := handler(ctx, req)
	if err != nil {
		logger.Warnw("RPC failed: ", "method", serverInfo.FullMethod, "request", req)
		return resp, err
	}

	logger.Infow("RPC executed: ", "method", serverInfo.FullMethod)
	return resp, nil
}

func valid(token []byte)
