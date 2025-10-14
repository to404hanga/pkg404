package loggerv2

import (
	"context"

	"github.com/to404hanga/pkg404/logger"
)

// NewConsoleLogger 创建输出到控制台的logger
// development: 是否为开发模式
func NewConsoleLogger(development bool) (*ZapContextLogger, error) {
	config := LoggerConfig{
		Output: OutputConfig{
			Type: OutputConsole,
		},
		Development: development,
	}
	return NewZapContextLoggerWithConfig(config)
}

// NewFileLogger 创建输出到文件的logger
// filePath: 文件路径
// development: 是否为开发模式
// autoCreate: 是否自动创建文件和目录
func NewFileLogger(filePath string, development bool, autoCreate bool) (*ZapContextLogger, error) {
	config := LoggerConfig{
		Output: OutputConfig{
			Type:           OutputFile,
			FilePath:       filePath,
			AutoCreateFile: autoCreate,
		},
		Development: development,
	}
	return NewZapContextLoggerWithConfig(config)
}

// NewBothLogger 创建同时输出到控制台和文件的logger
// filePath: 文件路径
// development: 是否为开发模式
// autoCreate: 是否自动创建文件和目录
func NewBothLogger(filePath string, development bool, autoCreate bool) (*ZapContextLogger, error) {
	config := LoggerConfig{
		Output: OutputConfig{
			Type:           OutputBoth,
			FilePath:       filePath,
			AutoCreateFile: autoCreate,
		},
		Development: development,
	}
	return NewZapContextLoggerWithConfig(config)
}

// NewDevelopmentLogger 创建开发模式的控制台logger
func NewDevelopmentLogger() (*ZapContextLogger, error) {
	return NewConsoleLogger(true)
}

// NewProductionLogger 创建生产模式的控制台logger
func NewProductionLogger() (*ZapContextLogger, error) {
	return NewConsoleLogger(false)
}

// NewDevelopmentFileLogger 创建开发模式的文件logger
func NewDevelopmentFileLogger(filePath string) (*ZapContextLogger, error) {
	return NewFileLogger(filePath, true, true)
}

// NewProductionFileLogger 创建生产模式的文件logger
func NewProductionFileLogger(filePath string) (*ZapContextLogger, error) {
	return NewFileLogger(filePath, false, true)
}

// MustNewConsoleLogger 创建控制台logger，失败时panic
func MustNewConsoleLogger(development bool) *ZapContextLogger {
	logger, err := NewConsoleLogger(development)
	if err != nil {
		panic(err)
	}
	return logger
}

// MustNewFileLogger 创建文件logger，失败时panic
func MustNewFileLogger(filePath string, development bool, autoCreate bool) *ZapContextLogger {
	logger, err := NewFileLogger(filePath, development, autoCreate)
	if err != nil {
		panic(err)
	}
	return logger
}

// MustNewBothLogger 创建双输出logger，失败时panic
func MustNewBothLogger(filePath string, development bool, autoCreate bool) *ZapContextLogger {
	logger, err := NewBothLogger(filePath, development, autoCreate)
	if err != nil {
		panic(err)
	}
	return logger
}

// 全局logger实例
var globalLogger Logger

// SetGlobalLogger 设置全局logger
func SetGlobalLogger(l Logger) {
	globalLogger = l
}

// GetGlobalLogger 获取全局logger
func GetGlobalLogger() Logger {
	if globalLogger == nil {
		// 默认创建一个开发模式的控制台logger
		globalLogger = MustNewConsoleLogger(true)
	}
	return globalLogger
}

// 全局便利函数
func Debug(msg string, args ...logger.Field) {
	GetGlobalLogger().Debug(msg, args...)
}

func Info(msg string, args ...logger.Field) {
	GetGlobalLogger().Info(msg, args...)
}

func Warn(msg string, args ...logger.Field) {
	GetGlobalLogger().Warn(msg, args...)
}

func Error(msg string, args ...logger.Field) {
	GetGlobalLogger().Error(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...logger.Field) {
	GetGlobalLogger().DebugContext(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...logger.Field) {
	GetGlobalLogger().InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...logger.Field) {
	GetGlobalLogger().WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...logger.Field) {
	GetGlobalLogger().ErrorContext(ctx, msg, args...)
}

func With(args ...logger.Field) Logger {
	return GetGlobalLogger().With(args...)
}

func WithContext(ctx context.Context) Logger {
	return GetGlobalLogger().WithContext(ctx)
}
