package slidingwindow

import (
	"context"
	"sync"
	"time"

	"github.com/to404hanga/pkg404/limiter"
	"github.com/to404hanga/pkg404/stl/queue"
)

type SlidingWindowLimiter struct {
	window    time.Duration
	queue     *queue.PriorityQueue[time.Time]
	lock      sync.Mutex
	threshold int
}

var _ limiter.Limiter = (*SlidingWindowLimiter)(nil)

func NewSlidingWindowLimiter(window time.Duration, threshold int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		window: window,
		queue: queue.NewPriorityQueueFunc[time.Time](func(left, right time.Time) bool {
			return left.Before(right)
		}),
		threshold: threshold,
	}
}

func (s *SlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now()

	// 快路径检测，原本就没满直接执行
	if s.queue.Len() < s.threshold {
		s.queue.Push(now)
		return true, nil
	}

	windowStart := time.Now().Add(-s.window)
	for {
		first := s.queue.Top()
		if first.Before(windowStart) {
			s.queue.Pop()
		} else {
			break
		}
	}
	if s.queue.Len() < s.threshold {
		s.queue.Push(now)
		return true, nil
	}
	return false, nil
}
