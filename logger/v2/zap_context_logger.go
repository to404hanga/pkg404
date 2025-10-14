package loggerv2

import (
	"context"

	"github.com/to404hanga/pkg404/logger"
	"go.uber.org/zap"
)

// ZapContextLogger 基于zap实现的支持context的logger
type ZapContextLogger struct {
	logger *zap.Logger
}

// 确保实现了所有接口
var _ Logger = (*ZapContextLogger)(nil)
var _ ContextLogger = (*ZapContextLogger)(nil)

// NewZapContextLogger 创建新的ZapContextLogger实例
func NewZapContextLogger(zapLogger *zap.Logger) *ZapContextLogger {
	return &ZapContextLogger{
		logger: zapLogger,
	}
}

// NewZapContextLoggerWithConfig 根据配置创建ZapContextLogger实例
func NewZapContextLoggerWithConfig(config LoggerConfig) (*ZapContextLogger, error) {
	zapLogger, err := createZapLogger(config)
	if err != nil {
		return nil, err
	}

	return &ZapContextLogger{
		logger: zapLogger,
	}, nil
}

// Debug 记录debug级别日志
func (l *ZapContextLogger) Debug(msg string, args ...logger.Field) {
	l.logger.Debug(msg, l.toZapFields(args)...)
}

// Info 记录info级别日志
func (l *ZapContextLogger) Info(msg string, args ...logger.Field) {
	l.logger.Info(msg, l.toZapFields(args)...)
}

// Warn 记录warn级别日志
func (l *ZapContextLogger) Warn(msg string, args ...logger.Field) {
	l.logger.Warn(msg, l.toZapFields(args)...)
}

// Error 记录error级别日志
func (l *ZapContextLogger) Error(msg string, args ...logger.Field) {
	l.logger.Error(msg, l.toZapFields(args)...)
}

// DebugContext 记录debug级别日志，支持context
func (l *ZapContextLogger) DebugContext(ctx context.Context, msg string, args ...logger.Field) {
	allFields := l.mergeContextFields(ctx, args)
	l.logger.Debug(msg, l.toZapFields(allFields)...)
}

// InfoContext 记录info级别日志，支持context
func (l *ZapContextLogger) InfoContext(ctx context.Context, msg string, args ...logger.Field) {
	allFields := l.mergeContextFields(ctx, args)
	l.logger.Info(msg, l.toZapFields(allFields)...)
}

// WarnContext 记录warn级别日志，支持context
func (l *ZapContextLogger) WarnContext(ctx context.Context, msg string, args ...logger.Field) {
	allFields := l.mergeContextFields(ctx, args)
	l.logger.Warn(msg, l.toZapFields(allFields)...)
}

// ErrorContext 记录error级别日志，支持context
func (l *ZapContextLogger) ErrorContext(ctx context.Context, msg string, args ...logger.Field) {
	allFields := l.mergeContextFields(ctx, args)
	l.logger.Error(msg, l.toZapFields(allFields)...)
}

// With 返回带有指定fields的新logger实例
func (l *ZapContextLogger) With(args ...logger.Field) Logger {
	// 直接在底层zap logger上添加fields，返回新的ZapContextLogger
	newZapLogger := l.logger.With(l.toZapFields(args)...)
	return &ZapContextLogger{
		logger: newZapLogger,
	}
}

// WithContext 返回带有context的新logger实例
func (l *ZapContextLogger) WithContext(ctx context.Context) Logger {
	contextFields := FieldsFromContext(ctx)

	// 直接在底层zap logger上添加context fields
	newZapLogger := l.logger.With(l.toZapFields(contextFields)...)
	return &ZapContextLogger{
		logger: newZapLogger,
	}
}

// ExtractFields 从context中提取所有fields
func (l *ZapContextLogger) ExtractFields(ctx context.Context) []logger.Field {
	return FieldsFromContext(ctx)
}

// WithFields 向context中添加fields
func (l *ZapContextLogger) WithFields(ctx context.Context, fields ...logger.Field) context.Context {
	return WithFieldsToContext(ctx, fields...)
}

// mergeContextFields 合并context中的fields和传入的fields
func (l *ZapContextLogger) mergeContextFields(ctx context.Context, args []logger.Field) []logger.Field {
	// 获取context中的fields
	contextFields := FieldsFromContext(ctx)

	// 合并context fields和传入的args
	allFields := make([]logger.Field, 0, len(contextFields)+len(args))
	allFields = append(allFields, contextFields...)
	allFields = append(allFields, args...)

	return allFields
}

// toZapFields 将logger.Field转换为zap.Field
func (l *ZapContextLogger) toZapFields(fields []logger.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Val))
	}
	return zapFields
}
