package counter

import (
	"context"
	"sync/atomic"

	"github.com/to404hanga/pkg404/limiter"
)

// CounterLimiter 基于计数器实现的限流器
type CounterLimiter struct {
	cnt       atomic.Int32
	threshold int32
}

var _ limiter.Limiter = (*CounterLimiter)(nil)

func NewCounterLimiter(threshold int32) *CounterLimiter {
	return &CounterLimiter{
		threshold: threshold,
		cnt:       atomic.Int32{},
	}
}

func (c *CounterLimiter) Limit(ctx context.Context, key string) (bool, error) {
	cnt := c.cnt.Add(1)
	defer func() {
		c.cnt.Add(-1)
	}()
	if cnt <= c.threshold {
		return true, nil
	}
	return false, nil
}
