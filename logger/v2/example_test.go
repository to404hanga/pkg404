package loggerv2

import (
	"context"
	"testing"

	"github.com/to404hanga/pkg404/logger"
	"go.uber.org/zap"
)

// TestExampleBasicUsage 展示基本用法
func TestExampleBasicUsage(t *testing.T) {
	// 创建zap logger
	zapLogger, _ := zap.NewDevelopment()
	l := NewZapContextLogger(zapLogger)

	// 基本日志记录
	l.Info("Hello world")
	l.Info("User login", logger.String("user_id", "123"), logger.String("username", "john"))

	// 使用With方法
	userLogger := l.With(logger.String("user_id", "123"))
	userLogger.Info("User action", logger.String("action", "click"))
}

// TestExampleContextUsage 展示context用法
func TestExampleContextUsage(t *testing.T) {
	zapLogger, _ := zap.NewDevelopment()
	l := NewZapContextLogger(zapLogger)

	// 方式1：使用WithFieldsToContext
	ctx := context.Background()
	ctx = WithFieldsToContext(ctx, logger.String("request_id", "req-123"))

	// 使用context记录日志
	l.InfoContext(ctx, "Processing request")

	// 添加更多fields到context
	ctx = WithFieldsToContext(ctx, logger.String("user_id", "user-456"))
	l.InfoContext(ctx, "User authenticated")

	// 方式2：使用更直观的ContextWithFields
	ctx2 := ContextWithFields(context.Background(),
		logger.String("session_id", "sess-789"),
		logger.String("ip", "192.168.1.1"))
	l.InfoContext(ctx2, "Session created")

	// 方式3：使用NewContextWithFields直接创建新context
	ctx3 := NewContextWithFields(
		logger.String("component", "auth"),
		logger.String("version", "v1.0"))
	l.InfoContext(ctx3, "Authentication module initialized")
}

// TestExampleFileOutput 展示文件输出功能
func TestExampleFileOutput(t *testing.T) {
	// 创建输出到文件的logger
	fileLogger, err := NewFileLogger("logs/app.log", true, true)
	if err != nil {
		t.Fatal(err)
	}

	// 记录日志到文件
	fileLogger.Info("Application started", logger.String("version", "1.0.0"))

	// 使用context记录日志
	ctx := NewContextWithFields(
		logger.String("request_id", "req-123"),
		logger.String("user_id", "user-456"))
	fileLogger.InfoContext(ctx, "User login successful")

	// 创建同时输出到控制台和文件的logger
	bothLogger, err := NewBothLogger("logs/app.log", true, true)
	if err != nil {
		panic(err)
	}

	bothLogger.Info("This message will appear in both console and file")
}

// TestExampleGlobalLogger 展示全局logger的使用
func TestExampleGlobalLogger(t *testing.T) {
	// 设置全局logger为文件logger
	fileLogger, _ := NewDevelopmentFileLogger("logs/global.log")
	SetGlobalLogger(fileLogger)

	// 使用全局函数记录日志
	Info("Global logger message", logger.String("component", "main"))

	// 使用context
	ctx := NewContextWithFields(logger.String("trace_id", "trace-789"))
	InfoContext(ctx, "Global context message")
}

// TestExampleAdapterUsage 展示适配器用法
func TestExampleAdapterUsage(t *testing.T) {
	// 使用原有的v1 logger
	zapLogger, _ := zap.NewDevelopment()
	v1Logger := logger.NewZapLogger(zapLogger)

	// 适配为v2 logger
	v2Logger := NewLoggerAdapter(v1Logger)

	// 现在可以使用v2的功能
	ctx := WithFieldsToContext(context.Background(), logger.String("session_id", "sess-789"))
	v2Logger.InfoContext(ctx, "Using adapted logger")
}

// TestLoggerV2Compatibility 测试v2 logger的兼容性
func TestLoggerV2Compatibility(t *testing.T) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	// 测试ZapContextLogger
	l := NewZapContextLogger(zapLogger)

	// 测试基本接口
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	// 测试context接口
	ctx := context.Background()
	ctx = WithFieldsToContext(ctx, logger.String("test", "value"))

	l.DebugContext(ctx, "debug with context")
	l.InfoContext(ctx, "info with context")
	l.WarnContext(ctx, "warn with context")
	l.ErrorContext(ctx, "error with context")

	// 测试With方法
	withLogger := l.With(logger.String("component", "test"))
	withLogger.Info("message with component")

	// 测试WithContext方法
	contextLogger := l.WithContext(ctx)
	contextLogger.Info("message with context logger")
}

// TestAdapter 测试适配器功能
func TestAdapter(t *testing.T) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	// 创建v1 logger
	v1Logger := logger.NewZapLogger(zapLogger)

	// 适配为v2 logger
	v2Logger := NewLoggerAdapter(v1Logger)

	// 测试基本功能
	v2Logger.Info("adapter test")

	// 测试context功能
	ctx := WithFieldsToContext(context.Background(), logger.String("adapter", "test"))
	v2Logger.InfoContext(ctx, "adapter context test")

	// 测试With方法
	withLogger := v2Logger.With(logger.String("with", "test"))
	withLogger.Info("adapter with test")
}

// TestFileOutput 测试文件输出功能
func TestFileOutput(t *testing.T) {
	// 测试文件logger
	fileLogger, err := NewFileLogger("test_logs/test.log", true, true)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	fileLogger.Info("Test file output", logger.String("test", "value"))

	// 测试双输出logger
	bothLogger, err := NewBothLogger("test_logs/both.log", true, true)
	if err != nil {
		t.Fatalf("Failed to create both logger: %v", err)
	}

	bothLogger.Info("Test both output", logger.String("test", "both"))

	// 测试context
	ctx := NewContextWithFields(logger.String("context", "test"))
	bothLogger.InfoContext(ctx, "Test context with file output")
}

// TestHelperFunctions 测试便利函数
func TestHelperFunctions(t *testing.T) {
	// 测试开发模式logger
	devLogger, err := NewDevelopmentLogger()
	if err != nil {
		t.Fatalf("Failed to create development logger: %v", err)
	}
	devLogger.Info("Development logger test")

	// 测试生产模式logger
	prodLogger, err := NewProductionLogger()
	if err != nil {
		t.Fatalf("Failed to create production logger: %v", err)
	}
	prodLogger.Info("Production logger test")

	// 测试文件logger便利函数
	devFileLogger, err := NewDevelopmentFileLogger("test_logs/dev.log")
	if err != nil {
		t.Fatalf("Failed to create development file logger: %v", err)
	}
	devFileLogger.Info("Development file logger test")

	prodFileLogger, err := NewProductionFileLogger("test_logs/prod.log")
	if err != nil {
		t.Fatalf("Failed to create production file logger: %v", err)
	}
	prodFileLogger.Info("Production file logger test")
}

// TestGlobalLogger 测试全局logger
func TestGlobalLogger(t *testing.T) {
	// 设置全局logger
	fileLogger, err := NewDevelopmentFileLogger("test_logs/global_test.log")
	if err != nil {
		t.Fatalf("Failed to create global logger: %v", err)
	}
	SetGlobalLogger(fileLogger)

	// 测试全局函数
	Info("Global info test", logger.String("global", "test"))
	Debug("Global debug test", logger.String("global", "debug"))

	// 测试context
	ctx := NewContextWithFields(logger.String("global_context", "test"))
	InfoContext(ctx, "Global context test")
}
