package ratelimit

import (
	"context"
	"strings"

	"github.com/to404hanga/pkg404/limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	limiter limiter.Limiter
	key     string
}

func NewInterceptorBuilder(limiter limiter.Limiter, key string) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key}
}

func (b *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil || limited {
			ctx = context.WithValue(ctx, "downgrade", "true")
		}
		return handler(ctx, req)
	}
}

func (b *InterceptorBuilder) BuildServerStreamInterceptorService() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if strings.HasPrefix(info.FullMethod, "/UserService") {
			limited, err := b.limiter.Limit(ctx, b.key)
			if err != nil {
				return nil, status.Errorf(codes.ResourceExhausted, "限流")
			}
			if limited {
				return nil, status.Errorf(codes.ResourceExhausted, "限流")
			}
		}
		return handler(ctx, req)
	}
}
