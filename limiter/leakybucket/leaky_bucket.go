package leakybucket

import (
	"context"
	"sync"
	"time"

	"github.com/to404hanga/pkg404/limiter"
)

type LeakyBucketLimiter struct {
	interval  time.Duration // 隔多久产生一个令牌
	closeCh   chan struct{}
	closeOnce sync.Once
}

func NewLeakyBucketLimiter(interval time.Duration) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		interval: interval,
		closeCh:  make(chan struct{}),
	}
}

var _ limiter.Limiter = (*LeakyBucketLimiter)(nil)

func (t *LeakyBucketLimiter) Limit(ctx context.Context, key string) (bool, error) {
	ticker := time.NewTicker(t.interval)

	select {
	case <-ticker.C:
		return true, nil
	case <-t.closeCh:
		return false, nil
	case <-ctx.Done(): // 等到超时了再退出
		return false, ctx.Err()
	}
}

func (t *LeakyBucketLimiter) Close() error {
	t.closeOnce.Do(func() {
		close(t.closeCh)
	})
	return nil
}
