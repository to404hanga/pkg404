package circuitbreaker

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
}

func (ib *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		err = ib.breaker.Allow()
		if err != nil {
			resp, err = handler(ctx, req)
			if err == nil {
				ib.breaker.MarkSuccess()
			} else {
				ib.breaker.MarkFailed()
			}
			return
		}
		ib.breaker.MarkFailed()
		return nil, status.Error(codes.Unavailable, "熔断")
	}
}
