package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestDoWithTimeout 测试context超时时的行为
func TestDoWithTimeout(t *testing.T) {
	// 创建一个500ms超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	start := time.Now()

	// 模拟一个总是失败的函数
	err := Do(ctx, func() error {
		return errors.New("always fail")
	}, RetryOptions{
		RetryTimes:        10,
		BaseInterval:      200 * time.Millisecond,
		BackoffMultiplier: 1.5,
	})

	elapsed := time.Since(start)

	// 应该因为context超时而退出
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}

	// 执行时间应该接近500ms，而不是等待所有重试完成
	if elapsed > 600*time.Millisecond {
		t.Errorf("Expected execution time around 500ms, got %v", elapsed)
	}

	t.Logf("Execution time: %v", elapsed)
}

// TestDoWithCancel 测试context取消时的行为
func TestDoWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	start := time.Now()

	// 在300ms后取消context
	go func() {
		time.Sleep(300 * time.Millisecond)
		cancel()
	}()

	// 模拟一个总是失败的函数
	err := Do(ctx, func() error {
		return errors.New("always fail")
	}, RetryOptions{
		RetryTimes:        10,
		BaseInterval:      200 * time.Millisecond,
		BackoffMultiplier: 1.5,
	})

	elapsed := time.Since(start)

	// 应该因为context被取消而退出
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}

	// 执行时间应该接近300ms
	if elapsed > 400*time.Millisecond {
		t.Errorf("Expected execution time around 300ms, got %v", elapsed)
	}

	t.Logf("Execution time: %v", elapsed)
}

// TestDoSuccess 测试成功情况
func TestDoSuccess(t *testing.T) {
	ctx := context.Background()

	callCount := 0
	err := Do(ctx, func() error {
		callCount++
		if callCount < 2 {
			return errors.New("fail")
		}
		return nil // 第二次调用成功
	}, RetryOptions{
		RetryTimes:        3,
		BaseInterval:      100 * time.Millisecond,
		BackoffMultiplier: 1.5,
	})

	if err != nil {
		t.Errorf("Expected success, got %v", err)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

// TestDoAllFailed 测试所有重试都失败的情况
func TestDoAllFailed(t *testing.T) {
	ctx := context.Background()

	callCount := 0
	err := Do(ctx, func() error {
		callCount++
		return errors.New("always fail")
	}, RetryOptions{
		RetryTimes:        3,
		BaseInterval:      50 * time.Millisecond,
		BackoffMultiplier: 1.5,
	})

	// 所有重试都失败时，应该返回nil（根据当前实现）
	if err != nil {
		t.Errorf("Expected nil when all retries failed, got %v", err)
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", callCount)
	}
}
