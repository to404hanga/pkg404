package loggerv2

import (
	"context"
	"testing"

	"github.com/to404hanga/pkg404/logger"
	"go.uber.org/zap"
)

// TestLoggerInterface 测试 Logger 接口的定义
func TestLoggerInterface(t *testing.T) {
	// 验证 ZapCtxLogger 实现了 Logger 接口
	var _ Logger = (*ZapCtxLogger)(nil)

	// 创建一个实际的实例来验证接口实现
	zapLogger := zap.NewNop()
	ctxLogger := NewZapCtxLogger(zapLogger)

	// 验证可以赋值给接口类型
	var loggerInterface Logger = ctxLogger
	if loggerInterface == nil {
		t.Error("ZapCtxLogger should implement Logger interface")
	}
}

// TestLoggerInterfaceBasicMethods 测试 Logger 接口的基本方法
func TestLoggerInterfaceBasicMethods(t *testing.T) {
	zapLogger := zap.NewNop()
	var logger Logger = NewZapCtxLogger(zapLogger)

	// 测试基本日志方法（这些方法来自嵌入的 logger.Logger 接口）
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// 如果没有 panic，说明方法调用成功
	t.Log("Basic logging methods work correctly")
}

// TestLoggerInterfaceContextMethods 测试 Logger 接口的 context 方法
func TestLoggerInterfaceContextMethods(t *testing.T) {
	zapLogger := zap.NewNop()
	var logger Logger = NewZapCtxLogger(zapLogger)

	ctx := context.Background()

	// 测试 context 日志方法
	logger.DebugContext(ctx, "debug with context")
	logger.InfoContext(ctx, "info with context")
	logger.WarnContext(ctx, "warn with context")
	logger.ErrorContext(ctx, "error with context")

	// 如果没有 panic，说明方法调用成功
	t.Log("Context logging methods work correctly")
}

// TestLoggerInterfaceChainMethods 测试 Logger 接口的链式调用方法
func TestLoggerInterfaceChainMethods(t *testing.T) {
	zapLogger := zap.NewNop()
	var l Logger = NewZapCtxLogger(zapLogger)

	ctx := context.Background()
	fields := []logger.Field{
		{Key: "test", Val: "value"},
	}

	// 测试链式调用方法
	newLogger1 := l.WithContext(ctx)
	if newLogger1 == nil {
		t.Error("WithContext should not return nil")
	}

	newLogger2 := l.WithFields(fields...)
	if newLogger2 == nil {
		t.Error("WithFields should not return nil")
	}

	// 验证返回的也是 Logger 接口类型
	var _ Logger = newLogger1
	var _ Logger = newLogger2

	t.Log("Chain methods work correctly")
}

// TestLoggerInterfaceCompatibility 测试 Logger 接口与父级 logger.Logger 的兼容性
func TestLoggerInterfaceCompatibility(t *testing.T) {
	zapLogger := zap.NewNop()
	ctxLogger := NewZapCtxLogger(zapLogger)

	// 验证可以赋值给父级接口
	var parentLogger logger.Logger = ctxLogger
	if parentLogger == nil {
		t.Error("ZapCtxLogger should be compatible with parent logger.Logger interface")
	}

	// 测试父级接口的方法
	parentLogger.Debug("debug message")
	parentLogger.Info("info message")
	parentLogger.Warn("warn message")
	parentLogger.Error("error message")

	t.Log("Parent interface compatibility works correctly")
}

// TestLoggerInterfaceMethodSignatures 测试 Logger 接口方法签名
func TestLoggerInterfaceMethodSignatures(t *testing.T) {
	zapLogger := zap.NewNop()
	var l Logger = NewZapCtxLogger(zapLogger)

	ctx := context.Background()
	field := logger.Field{Key: "test", Val: "value"}

	// 测试方法签名是否正确

	// Context 方法应该接受 context 和可变参数
	l.DebugContext(ctx, "message")
	l.DebugContext(ctx, "message", field)
	l.DebugContext(ctx, "message", field, field)

	l.InfoContext(ctx, "message")
	l.InfoContext(ctx, "message", field)

	l.WarnContext(ctx, "message")
	l.WarnContext(ctx, "message", field)

	l.ErrorContext(ctx, "message")
	l.ErrorContext(ctx, "message", field)

	// WithContext 应该返回 Logger 接口
	newLogger := l.WithContext(ctx)
	var _ Logger = newLogger

	// WithFields 应该接受可变参数并返回 Logger 接口
	newLogger2 := l.WithFields()
	newLogger3 := l.WithFields(field)
	newLogger4 := l.WithFields(field, field)

	var _ Logger = newLogger2
	var _ Logger = newLogger3
	var _ Logger = newLogger4

	t.Log("Method signatures are correct")
}

// TestLoggerInterfaceNilSafety 测试 Logger 接口的 nil 安全性
func TestLoggerInterfaceNilSafety(t *testing.T) {
	zapLogger := zap.NewNop()
	var l Logger = NewZapCtxLogger(zapLogger)

	// 测试传入 nil context
	l.DebugContext(nil, "message with nil context")
	l.InfoContext(nil, "message with nil context")
	l.WarnContext(nil, "message with nil context")
	l.ErrorContext(nil, "message with nil context")

	// 测试 WithContext 传入 nil
	newLogger := l.WithContext(nil)
	if newLogger == nil {
		t.Error("WithContext should not return nil even with nil context")
	}

	// 测试 WithFields 传入 nil 或空参数
	newLogger2 := l.WithFields()
	if newLogger2 == nil {
		t.Error("WithFields should not return nil with empty arguments")
	}

	t.Log("Nil safety tests passed")
}

// TestLoggerInterfaceFieldTypes 测试 Logger 接口处理不同类型的字段
func TestLoggerInterfaceFieldTypes(t *testing.T) {
	zapLogger := zap.NewNop()
	var l Logger = NewZapCtxLogger(zapLogger)

	ctx := context.Background()

	// 测试不同类型的字段值
	fields := []logger.Field{
		{Key: "string", Val: "string_value"},
		{Key: "int", Val: 42},
		{Key: "bool", Val: true},
		{Key: "float", Val: 3.14},
		{Key: "slice", Val: []string{"a", "b", "c"}},
		{Key: "map", Val: map[string]int{"x": 1}},
		{Key: "nil", Val: nil},
	}

	// 测试所有日志级别都能处理这些字段类型
	l.DebugContext(ctx, "debug with various field types", fields...)
	l.InfoContext(ctx, "info with various field types", fields...)
	l.WarnContext(ctx, "warn with various field types", fields...)
	l.ErrorContext(ctx, "error with various field types", fields...)

	// 测试 WithFields 也能处理这些字段类型
	newLogger := l.WithFields(fields...)
	if newLogger == nil {
		t.Error("WithFields should handle various field types")
	}

	t.Log("Field type handling tests passed")
}

// TestLoggerInterfaceChaining 测试 Logger 接口的链式调用
func TestLoggerInterfaceChaining(t *testing.T) {
	zapLogger := zap.NewNop()
	var l Logger = NewZapCtxLogger(zapLogger)

	ctx := context.Background()
	field1 := logger.Field{Key: "field1", Val: "value1"}
	field2 := logger.Field{Key: "field2", Val: "value2"}

	// 测试复杂的链式调用
	l.WithContext(ctx).
		WithFields(field1).
		WithFields(field2).
		InfoContext(context.Background(), "chained call test")

	// 测试链式调用的返回类型都是 Logger 接口
	step1 := l.WithContext(ctx)
	var _ Logger = step1

	step2 := step1.WithFields(field1)
	var _ Logger = step2

	step3 := step2.WithFields(field2)
	var _ Logger = step3

	t.Log("Interface chaining tests passed")
}

// MockLogger 用于测试的 mock 实现
type MockLogger struct {
	debugCalls   int
	infoCalls    int
	warnCalls    int
	errorCalls   int
	contextCalls int
	withCalls    int
}

func (m *MockLogger) Debug(msg string, args ...logger.Field) {
	m.debugCalls++
}

func (m *MockLogger) Info(msg string, args ...logger.Field) {
	m.infoCalls++
}

func (m *MockLogger) Warn(msg string, args ...logger.Field) {
	m.warnCalls++
}

func (m *MockLogger) Error(msg string, args ...logger.Field) {
	m.errorCalls++
}

func (m *MockLogger) DebugContext(ctx context.Context, msg string, args ...logger.Field) {
	m.debugCalls++
	m.contextCalls++
}

func (m *MockLogger) InfoContext(ctx context.Context, msg string, args ...logger.Field) {
	m.infoCalls++
	m.contextCalls++
}

func (m *MockLogger) WarnContext(ctx context.Context, msg string, args ...logger.Field) {
	m.warnCalls++
	m.contextCalls++
}

func (m *MockLogger) ErrorContext(ctx context.Context, msg string, args ...logger.Field) {
	m.errorCalls++
	m.contextCalls++
}

func (m *MockLogger) WithContext(ctx context.Context) Logger {
	m.withCalls++
	return m
}

func (m *MockLogger) WithFields(args ...logger.Field) Logger {
	m.withCalls++
	return m
}

// TestLoggerInterfaceImplementation 测试自定义 Logger 实现
func TestLoggerInterfaceImplementation(t *testing.T) {
	var l Logger = &MockLogger{}

	ctx := context.Background()
	field := logger.Field{Key: "test", Val: "value"}

	// 测试所有方法都能正常调用
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")

	l.DebugContext(ctx, "debug context")
	l.InfoContext(ctx, "info context")
	l.WarnContext(ctx, "warn context")
	l.ErrorContext(ctx, "error context")

	l.WithContext(ctx)
	l.WithFields(field)

	// 验证调用计数
	mockLogger := l.(*MockLogger)
	if mockLogger.debugCalls != 2 {
		t.Errorf("expected 2 debug calls, got %d", mockLogger.debugCalls)
	}
	if mockLogger.infoCalls != 2 {
		t.Errorf("expected 2 info calls, got %d", mockLogger.infoCalls)
	}
	if mockLogger.warnCalls != 2 {
		t.Errorf("expected 2 warn calls, got %d", mockLogger.warnCalls)
	}
	if mockLogger.errorCalls != 2 {
		t.Errorf("expected 2 error calls, got %d", mockLogger.errorCalls)
	}
	if mockLogger.contextCalls != 4 {
		t.Errorf("expected 4 context calls, got %d", mockLogger.contextCalls)
	}
	if mockLogger.withCalls != 2 {
		t.Errorf("expected 2 with calls, got %d", mockLogger.withCalls)
	}

	t.Log("Custom Logger implementation tests passed")
}
