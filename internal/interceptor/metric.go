package interceptor

import (
	"context"

	"github.com/pavel-pj/auth/internal/metric"
	"google.golang.org/grpc"
)

func MetricInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	metric.IncRequestCounter()
	return handler(ctx, req)
}
