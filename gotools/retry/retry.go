package retry

import (
	"context"
	"time"
)

const (
	DefaultRetryTimes        = 3
	DefaultBaseInterval      = 100 * time.Millisecond
	DefaultBackoffMultiplier = 1.5
)

type RetryOptions struct {
	RetryTimes        int
	BaseInterval      time.Duration
	BackoffMultiplier float64 // 退避倍率，每次重试间隔会乘以这个倍率
}

// Do 执行带重试的函数，支持context超时控制
func Do(ctx context.Context, fn func() error, opts ...RetryOptions) error {
	retryTimes := DefaultRetryTimes
	interval := DefaultBaseInterval
	backoffMultiplier := DefaultBackoffMultiplier
	if len(opts) > 0 {
		retryTimes = opts[0].RetryTimes
		interval = opts[0].BaseInterval
		backoffMultiplier = opts[0].BackoffMultiplier
	}

	for i := 0; i < retryTimes; i++ {
		// 在每次重试前检查context状态
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := fn(); err == nil {
			return nil
		}

		// 如果不是最后一次重试，则等待后重试
		if i < retryTimes-1 {
			// 使用可被context中断的sleep
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
				// 继续下一次重试
			}
			// 每次重试后，更新等待时间
			interval = time.Duration(float64(interval) * backoffMultiplier)
		}
	}
	return nil
}
