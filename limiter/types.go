package limiter

import "context"

//go:generate mockgen -source=./types.go -package=limitermocks -destination=./mocks/limiter.mock.go Limiter
type Limiter interface {
	// Limit 是否触发限流
	Limit(ctx context.Context, key string) (bool, error)
}
