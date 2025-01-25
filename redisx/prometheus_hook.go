package redisx

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type PrometheusHook struct {
	vector *prometheus.SummaryVec
}

var _ redis.Hook = (*PrometheusHook)(nil)

func NewPrometheusHook(opt prometheus.SummaryOpts) redis.Hook {
	return &PrometheusHook{
		vector: prometheus.NewSummaryVec(opt, []string{"cmd", "key_exists"}),
	}
}

func (h *PrometheusHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (h *PrometheusHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		var err error
		defer func() {
			duration := time.Since(start).Milliseconds()
			keyExists := err == redis.Nil
			h.vector.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExists)).Observe(float64(duration))
		}()
		err = next(ctx, cmd)
		return err
	}
}

func (h *PrometheusHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
