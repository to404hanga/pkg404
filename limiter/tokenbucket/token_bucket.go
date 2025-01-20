package tokenbucket

import (
	"context"
	"sync"
	"time"

	"github.com/to404hanga/pkg404/limiter"
)

type TokenBucket struct {
	interval  time.Duration // 隔多久产生一个令牌
	buckets   chan struct{}
	closeCh   chan struct{}
	closeOnce sync.Once
}

func NewTokenBucket(interval time.Duration, bucketSize int) *TokenBucket {
	return &TokenBucket{
		interval: interval,
		buckets:  make(chan struct{}, bucketSize),
		closeCh:  make(chan struct{}),
	}
}

var _ limiter.Limiter = (*TokenBucket)(nil)

func (t *TokenBucket) Limit(ctx context.Context, key string) (bool, error) {
	ticker := time.NewTicker(t.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				select {
				case t.buckets <- struct{}{}:
				default:
					// 令牌桶满了
				}
			case <-t.closeCh:
				return
			}
		}
	}()

	select {
	case <-t.buckets:
		return true, nil
	case <-ctx.Done(): // 等到超时了再退出
		return false, ctx.Err()
	}
}

func (t *TokenBucket) Close() error {
	t.closeOnce.Do(func() {
		close(t.closeCh)
	})
	return nil
}
