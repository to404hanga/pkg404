package loggerv2

import (
	"context"
	"io"

	"github.com/to404hanga/pkg404/logger"
)

// Logger 定义了支持context的日志接口
// 兼容原有logger.Logger接口，同时支持context传递
type Logger interface {
	// 兼容原有接口
	logger.Logger

	// 支持context的新接口
	DebugContext(ctx context.Context, msg string, args ...logger.Field)
	InfoContext(ctx context.Context, msg string, args ...logger.Field)
	WarnContext(ctx context.Context, msg string, args ...logger.Field)
	ErrorContext(ctx context.Context, msg string, args ...logger.Field)

	// 支持链式调用的With方法
	With(args ...logger.Field) Logger
	WithContext(ctx context.Context) Logger
}

// ContextLogger 扩展接口，支持从context中提取fields
type ContextLogger interface {
	Logger

	// 从context中提取所有fields
	ExtractFields(ctx context.Context) []logger.Field

	// 向context中添加fields
	WithFields(ctx context.Context, fields ...logger.Field) context.Context
}

// contextKey 用于在context中存储logger fields的key类型
type contextKey string

const (
	// FieldsKey context中存储fields的key
	FieldsKey contextKey = "logger_fields"
)

// FieldsFromContext 从context中提取fields
func FieldsFromContext(ctx context.Context) []logger.Field {
	if ctx == nil {
		return nil
	}

	fields, ok := ctx.Value(FieldsKey).([]logger.Field)
	if !ok {
		return nil
	}

	return fields
}

// WithFieldsToContext 向context中添加fields
func WithFieldsToContext(ctx context.Context, fields ...logger.Field) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	existingFields := FieldsFromContext(ctx)
	allFields := append(existingFields, fields...)

	return context.WithValue(ctx, FieldsKey, allFields)
}

// ContextWithFields 创建携带指定fields的context（外部可用的便利函数）
// 这是一个更直观的函数名，用于外部调用
func ContextWithFields(ctx context.Context, fields ...logger.Field) context.Context {
	return WithFieldsToContext(ctx, fields...)
}

// NewContextWithFields 基于background context创建携带指定fields的新context
func NewContextWithFields(fields ...logger.Field) context.Context {
	return WithFieldsToContext(context.Background(), fields...)
}

// OutputType 定义日志输出类型
type OutputType int

const (
	// OutputConsole 输出到控制台
	OutputConsole OutputType = iota
	// OutputFile 输出到文件
	OutputFile
	// OutputBoth 同时输出到控制台和文件
	OutputBoth
)

// OutputConfig 日志输出配置
type OutputConfig struct {
	// Type 输出类型
	Type OutputType
	// FilePath 文件路径（当Type为OutputFile或OutputBoth时使用）
	FilePath string
	// AutoCreateFile 是否自动创建文件和目录
	AutoCreateFile bool
	// Writer 自定义输出writer（可选）
	Writer io.Writer
}

// LoggerConfig 日志器配置
type LoggerConfig struct {
	// Output 输出配置
	Output OutputConfig
	// Development 是否为开发模式
	Development bool
	// Level 日志级别
	Level string
}
