package interceptor

import (
	"context"
	"time"

	"github.com/pavel-pj/auth/internal/metric"
	"google.golang.org/grpc"
)

func MetricInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	metric.IncRequestCounter()
	timestart := time.Now()

	res, err := handler(ctx, req)
	diffTime := time.Since(timestart)
	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
		metric.HistogramResponseTimeObserve("error", diffTime.Seconds())
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
		metric.HistogramResponseTimeObserve("success", diffTime.Seconds())
	}
	return res, err
}
