package interceptor

import (
	"context"
	"time"

	"github.com/pavel-pj/auth/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method ", info.FullMethod), zap.Any("req", req))
		return res, err
	}
	logger.Info("request", zap.String("method ", info.FullMethod), zap.Any("req", req), zap.Any("time", now))

	return res, err

}
