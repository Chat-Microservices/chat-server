package interceptor

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("req", req))
	}

	logger.Info(
		"request success",
		zap.String("method", info.FullMethod),
		zap.Any("req", req),
		zap.Any("res", res),
		zap.Duration("duration", time.Since(now)),
	)

	return res, err
}
