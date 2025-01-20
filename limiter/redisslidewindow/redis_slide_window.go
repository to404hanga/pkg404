package redisslidewindow

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/to404hanga/pkg404/limiter"
)

//go:embed slide_window.lua
var luaScript string

// RedisSlidingWindowLimiter 基于 Redis 的滑动窗口限流器
type RedisSlidingWindowLimiter struct {
	cmd      redis.Cmdable
	interval time.Duration
	rate     int // 阈值
}

var _ limiter.Limiter = (*RedisSlidingWindowLimiter)(nil)

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaScript, []string{key}, r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
