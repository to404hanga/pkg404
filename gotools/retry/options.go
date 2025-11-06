package retry

import "time"

type options struct {
	retryTimes        int
	baseInterval      time.Duration
	backoffMultiplier float64
	async             bool
	callback          func(error)
}

type Option func(*options)

func WithRetryTimes(times int) Option {
	return func(o *options) {
		if times > 0 {
			o.retryTimes = times
		}
	}
}

func WithBaseInterval(interval time.Duration) Option {
	return func(o *options) {
		if interval > 0 {
			o.baseInterval = interval
		}
	}
}

func WithBackoffMultiplier(multiplier float64) Option {
	return func(o *options) {
		if multiplier > 0 {
			o.backoffMultiplier = multiplier
		}
	}
}

func WithAsync(async bool) Option {
	return func(o *options) {
		o.async = async
	}
}

func WithCallback(callback func(error)) Option {
	return func(o *options) {
		o.callback = callback
	}
}
