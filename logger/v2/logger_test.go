package loggerv2

import (
	"context"
	"testing"

	"github.com/to404hanga/pkg404/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// createTestLogger 创建用于测试的 logger 和观察器
func createTestLogger() (*ZapCtxLogger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)
	ctxLogger := NewZapCtxLogger(zapLogger)
	return ctxLogger, logs
}

// TestNewZapCtxLogger 测试创建新的 ZapCtxLogger 实例
func TestNewZapCtxLogger(t *testing.T) {
	zapLogger := zap.NewNop()

	// 测试使用默认配置
	logger1 := NewZapCtxLogger(zapLogger)
	if logger1 == nil {
		t.Fatal("NewZapCtxLogger should not return nil")
	}
	if logger1.config == nil {
		t.Fatal("config should not be nil")
	}
	if !logger1.config.EnableContextExtraction {
		t.Error("default config should enable context extraction")
	}

	// 测试使用自定义配置
	customConfig := &LoggerConfig{
		EnableContextExtraction: false,
		ContextExtractor:        nil,
		ExtractAllContextValues: true,
	}
	logger2 := NewZapCtxLogger(zapLogger, customConfig)
	if logger2.config.EnableContextExtraction {
		t.Error("custom config should disable context extraction")
	}
	if logger2.config.ExtractAllContextValues != true {
		t.Error("custom config should set ExtractAllContextValues to true")
	}
}

// TestWithContext 测试 WithContext 方法
func TestWithContext(t *testing.T) {
	logger, _ := createTestLogger()
	ctx := context.WithValue(context.Background(), "test_key", "test_value")

	newLogger := logger.WithContext(ctx)
	if newLogger == nil {
		t.Fatal("WithContext should not return nil")
	}

	// 验证返回的是新实例
	if newLogger == logger {
		t.Error("WithContext should return a new instance")
	}

	// 验证 context 被正确设置
	ctxLogger := newLogger.(*ZapCtxLogger)
	if ctxLogger.ctx != ctx {
		t.Error("context should be set correctly")
	}
}

// TestWithFields 测试 WithFields 方法
func TestWithFields(t *testing.T) {
	log, _ := createTestLogger()
	fields := []logger.Field{
		{Key: "field1", Val: "value1"},
		{Key: "field2", Val: 42},
	}

	newLogger := log.WithFields(fields...)
	if newLogger == nil {
		t.Fatal("WithFields should not return nil")
	}

	// 验证返回的是新实例
	if newLogger == log {
		t.Error("WithFields should return a new instance")
	}

	// 验证字段被正确设置
	ctxLogger := newLogger.(*ZapCtxLogger)
	if len(ctxLogger.withFields) != len(fields) {
		t.Errorf("expected %d fields, got %d", len(fields), len(ctxLogger.withFields))
	}

	for i, field := range fields {
		if ctxLogger.withFields[i].Key != field.Key {
			t.Errorf("field %d key mismatch: expected %s, got %s", i, field.Key, ctxLogger.withFields[i].Key)
		}
		if ctxLogger.withFields[i].Val != field.Val {
			t.Errorf("field %d value mismatch: expected %v, got %v", i, field.Val, ctxLogger.withFields[i].Val)
		}
	}
}

// TestBasicLoggingMethods 测试基本的日志记录方法
func TestBasicLoggingMethods(t *testing.T) {
	l, logs := createTestLogger()

	// 测试 Debug
	l.Debug("debug message", logger.Field{Key: "level", Val: "debug"})
	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}
	entry := logs.All()[0]
	if entry.Level != zapcore.DebugLevel {
		t.Errorf("expected debug level, got %v", entry.Level)
	}
	if entry.Message != "debug message" {
		t.Errorf("expected 'debug message', got '%s'", entry.Message)
	}

	// 测试 Info
	logs.TakeAll() // 清空日志
	l.Info("info message", logger.Field{Key: "level", Val: "info"})
	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}
	entry = logs.All()[0]
	if entry.Level != zapcore.InfoLevel {
		t.Errorf("expected info level, got %v", entry.Level)
	}

	// 测试 Warn
	logs.TakeAll()
	l.Warn("warn message", logger.Field{Key: "level", Val: "warn"})
	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}
	entry = logs.All()[0]
	if entry.Level != zapcore.WarnLevel {
		t.Errorf("expected warn level, got %v", entry.Level)
	}

	// 测试 Error
	logs.TakeAll()
	l.Error("error message", logger.Field{Key: "level", Val: "error"})
	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}
	entry = logs.All()[0]
	if entry.Level != zapcore.ErrorLevel {
		t.Errorf("expected error level, got %v", entry.Level)
	}
}

// TestContextLoggingMethods 测试带 context 的日志记录方法
func TestContextLoggingMethods(t *testing.T) {
	l, logs := createTestLogger()
	ctx := context.WithValue(context.Background(), "user_id", "12345")

	// 测试 DebugContext
	l.DebugContext(ctx, "debug with context")
	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}
	entry := logs.All()[0]
	if entry.Level != zapcore.DebugLevel {
		t.Errorf("expected debug level, got %v", entry.Level)
	}

	// 测试 InfoContext
	logs.TakeAll()
	l.InfoContext(ctx, "info with context")
	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}
	entry = logs.All()[0]
	if entry.Level != zapcore.InfoLevel {
		t.Errorf("expected info level, got %v", entry.Level)
	}

	// 测试 WarnContext
	logs.TakeAll()
	l.WarnContext(ctx, "warn with context")
	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}
	entry = logs.All()[0]
	if entry.Level != zapcore.WarnLevel {
		t.Errorf("expected warn level, got %v", entry.Level)
	}

	// 测试 ErrorContext
	logs.TakeAll()
	l.ErrorContext(ctx, "error with context")
	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}
	entry = logs.All()[0]
	if entry.Level != zapcore.ErrorLevel {
		t.Errorf("expected error level, got %v", entry.Level)
	}
}

// TestContextExtraction 测试 context 值提取功能
func TestContextExtraction(t *testing.T) {
	l, logs := createTestLogger()

	// 创建包含多个值的 context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", "12345")
	ctx = context.WithValue(ctx, "request_id", "req-67890")

	l.InfoContext(ctx, "test message")

	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}

	entry := logs.All()[0]

	// 检查是否包含 context 中的值
	hasUserID := false
	hasRequestID := false

	for _, field := range entry.Context {
		if field.Key == "user_id" && field.String == "12345" {
			hasUserID = true
		}
		if field.Key == "request_id" && field.String == "req-67890" {
			hasRequestID = true
		}
	}

	if !hasUserID {
		t.Error("expected user_id to be extracted from context")
	}
	if !hasRequestID {
		t.Error("expected request_id to be extracted from context")
	}
}

// TestWithMethod 测试 With 方法的向后兼容性
func TestWithMethod(t *testing.T) {
	l, _ := createTestLogger()
	fields := []logger.Field{
		{Key: "service", Val: "test-service"},
		{Key: "version", Val: "1.0.0"},
	}

	newLogger := l.With(fields...)
	if newLogger == nil {
		t.Fatal("With should not return nil")
	}

	// 验证字段被正确设置
	if len(newLogger.withFields) != len(fields) {
		t.Errorf("expected %d fields, got %d", len(fields), len(newLogger.withFields))
	}
}

// TestChainedCalls 测试链式调用
func TestChainedCalls(t *testing.T) {
	l, logs := createTestLogger()
	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")

	// 测试链式调用
	l.WithContext(ctx).
		WithFields(logger.Field{Key: "component", Val: "test"}).
		InfoContext(context.WithValue(context.Background(), "span_id", "span-456"), "chained call test")

	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}

	entry := logs.All()[0]
	if entry.Message != "chained call test" {
		t.Errorf("expected 'chained call test', got '%s'", entry.Message)
	}

	// 检查字段是否正确包含
	hasComponent := false
	hasSpanID := false

	for _, field := range entry.Context {
		if field.Key == "component" && field.String == "test" {
			hasComponent = true
		}
		if field.Key == "span_id" && field.String == "span-456" {
			hasSpanID = true
		}
	}

	if !hasComponent {
		t.Error("expected component field to be present")
	}
	if !hasSpanID {
		t.Error("expected span_id to be extracted from context")
	}
}

// TestDisabledContextExtraction 测试禁用 context 提取
func TestDisabledContextExtraction(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)

	config := &LoggerConfig{
		EnableContextExtraction: false,
		ContextExtractor:        nil,
	}

	logger := NewZapCtxLogger(zapLogger, config)
	ctx := context.WithValue(context.Background(), "user_id", "12345")

	logger.InfoContext(ctx, "test message")

	if logs.Len() != 1 {
		t.Errorf("expected 1 log entry, got %d", logs.Len())
	}

	entry := logs.All()[0]

	// 检查不应该包含 context 中的值
	for _, field := range entry.Context {
		if field.Key == "user_id" {
			t.Error("user_id should not be extracted when context extraction is disabled")
		}
	}
}
