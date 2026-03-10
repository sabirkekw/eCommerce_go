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

	// validating JWT token
	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return nil, status.Errorf(codes.Unauthenticated, "no token")
	// }
	// token := md["authorization"]
	// isValid, err := valid(token, ctx.Value("jwtSecret").(string))
	// if err != nil {
	// 	if errors.Is(err, ErrTokenExpired) {
	// 		return nil, status.Errorf(codes.Unauthenticated, "token expired")
	// 	}
	// 	return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	// }
	// if !isValid {
	// 	return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	// }
	// logger.Debugw("Token valid", "method", serverInfo.FullMethod)

	// handling RPC
	resp, err := handler(ctx, req)
	if err != nil {
		logger.Warnw("RPC failed: ", "method", serverInfo.FullMethod, "request", req)
		return resp, err
	}

	logger.Debugw("RPC executed: ", "method", serverInfo.FullMethod)
	return resp, nil
}

// func valid(tokenSliced []string, jwtSecret string) (bool, error) {
// 	token := strings.Join(tokenSliced, "")
// 	token, _ = strings.CutPrefix(token, "Bearer ")
// 	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(jwtSecret), nil
// 	})

// 	if err != nil {
// 		if errors.Is(err, jwt.ErrTokenExpired) {
// 			return false, ErrTokenExpired
// 		}
// 		return false, ErrInvalidToken
// 	}

// 	if !parsedToken.Valid {
// 		return false, ErrInvalidToken
// 	}
// 	return true, nil
// }
