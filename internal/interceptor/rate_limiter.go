package interceptor

import (
	"context"

	ratelimiter "github.com/pavel-pj/auth/internal/rate_limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimiter struct {
	ratelimiter *ratelimiter.TokenBucketLimiter
}

func NewInterceptorRateLimiter(limiter *ratelimiter.TokenBucketLimiter) *RateLimiter {
	return &RateLimiter{
		ratelimiter: limiter,
	}
}

func (r *RateLimiter) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	if !r.ratelimiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "too many requests")
	}

	return handler(ctx, req)

}
