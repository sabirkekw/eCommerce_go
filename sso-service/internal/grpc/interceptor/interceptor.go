package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// getting logger from context
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		panic("failed to recieve logger from context")
	}
	logger.Infow("Recieved request: ", "method", serverInfo.FullMethod, "request", req)

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Infow("RPC failed: ", "method", serverInfo.FullMethod, "request", req)
		return resp, err
	}

	logger.Infow("RPC executed: ", "method", serverInfo.FullMethod)
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
