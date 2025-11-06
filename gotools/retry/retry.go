package retry

import (
	"context"
	"errors"
	"time"
)

func Do(ctx context.Context, fn func() error, opts ...Option) error {
	optsContainer := &options{
		retryTimes:        3,
		baseInterval:      100 * time.Millisecond,
		backoffMultiplier: 1.5,
		async:             false,
		callback:          nil,
	}

	for _, opt := range opts {
		opt(optsContainer)
	}

	if optsContainer.async {
		go func() {
			finalErr := doInternal(ctx, fn, optsContainer)
			if optsContainer.callback != nil {
				optsContainer.callback(finalErr)
			}
		}()
		return nil
	}

	return doInternal(ctx, fn, optsContainer)
}

func doInternal(ctx context.Context, fn func() error, opts *options) error {
	var err error
	interval := opts.baseInterval

	for i := 0; i < opts.retryTimes; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = fn()
		if err == nil {
			return nil // Success
		}

		if i < opts.retryTimes-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
			}
			interval = time.Duration(float64(interval) * opts.backoffMultiplier)
		}
	}

	return errors.New("function failed after all retries, last error: " + err.Error())
}
