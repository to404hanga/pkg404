package fixedwindow

import (
	"context"
	"sync"
	"time"

	"github.com/to404hanga/pkg404/limiter"
)

// FixedWindowLimiter 基于固定窗口实现的限流器
type FixedWindowLimiter struct {
	window          time.Duration
	lastWindowStart time.Time
	cnt             int
	threshold       int
	lock            sync.Mutex
}

var _ limiter.Limiter = (*FixedWindowLimiter)(nil)

func NewFixedWindowLimiter(window time.Duration, threshold int) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		window:          window,
		lastWindowStart: time.Now(),
		cnt:             0,
		threshold:       threshold,
	}
}

func (f *FixedWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	f.lock.Lock()
	now := time.Now()
	if now.After(f.lastWindowStart.Add(f.window)) {
		f.cnt = 0
		f.lastWindowStart = now
	}
	cnt := f.cnt + 1
	f.lock.Unlock()
	if cnt <= f.threshold {
		return true, nil
	}
	return false, nil
}
