package loggerv2

import (
	"context"

	"github.com/to404hanga/pkg404/logger"
)

// LoggerAdapter 将原有的logger.Logger适配为v2.Logger
type LoggerAdapter struct {
	logger logger.Logger
}

// 确保实现了Logger接口
var _ Logger = (*LoggerAdapter)(nil)

// NewLoggerAdapter 创建logger适配器
func NewLoggerAdapter(l logger.Logger) *LoggerAdapter {
	return &LoggerAdapter{logger: l}
}

// Debug 记录debug级别日志
func (a *LoggerAdapter) Debug(msg string, args ...logger.Field) {
	a.logger.Debug(msg, args...)
}

// Info 记录info级别日志
func (a *LoggerAdapter) Info(msg string, args ...logger.Field) {
	a.logger.Info(msg, args...)
}

// Warn 记录warn级别日志
func (a *LoggerAdapter) Warn(msg string, args ...logger.Field) {
	a.logger.Warn(msg, args...)
}

// Error 记录error级别日志
func (a *LoggerAdapter) Error(msg string, args ...logger.Field) {
	a.logger.Error(msg, args...)
}

// DebugContext 记录debug级别日志，忽略context
func (a *LoggerAdapter) DebugContext(ctx context.Context, msg string, args ...logger.Field) {
	// 从context中提取fields并合并
	contextFields := FieldsFromContext(ctx)
	allFields := append(contextFields, args...)
	a.logger.Debug(msg, allFields...)
}

// InfoContext 记录info级别日志，忽略context
func (a *LoggerAdapter) InfoContext(ctx context.Context, msg string, args ...logger.Field) {
	// 从context中提取fields并合并
	contextFields := FieldsFromContext(ctx)
	allFields := append(contextFields, args...)
	a.logger.Info(msg, allFields...)
}

// WarnContext 记录warn级别日志，忽略context
func (a *LoggerAdapter) WarnContext(ctx context.Context, msg string, args ...logger.Field) {
	// 从context中提取fields并合并
	contextFields := FieldsFromContext(ctx)
	allFields := append(contextFields, args...)
	a.logger.Warn(msg, allFields...)
}

// ErrorContext 记录error级别日志，忽略context
func (a *LoggerAdapter) ErrorContext(ctx context.Context, msg string, args ...logger.Field) {
	// 从context中提取fields并合并
	contextFields := FieldsFromContext(ctx)
	allFields := append(contextFields, args...)
	a.logger.Error(msg, allFields...)
}

// With 返回带有指定fields的新适配器
func (a *LoggerAdapter) With(args ...logger.Field) Logger {
	// 如果原logger支持With方法，使用它
	if withLogger, ok := a.logger.(*logger.ZapLogger); ok {
		return NewLoggerAdapter(withLogger.With(args...))
	}

	// 否则创建一个包装器
	return &LoggerAdapterWithFields{
		adapter: a,
		fields:  args,
	}
}

// WithContext 返回带有context的新适配器
func (a *LoggerAdapter) WithContext(ctx context.Context) Logger {
	contextFields := FieldsFromContext(ctx)
	if len(contextFields) == 0 {
		return a
	}

	return &LoggerAdapterWithFields{
		adapter: a,
		fields:  contextFields,
	}
}

// LoggerAdapterWithFields 带有预设fields的适配器
type LoggerAdapterWithFields struct {
	adapter *LoggerAdapter
	fields  []logger.Field
}

// 确保实现了Logger接口
var _ Logger = (*LoggerAdapterWithFields)(nil)

// Debug 记录debug级别日志
func (a *LoggerAdapterWithFields) Debug(msg string, args ...logger.Field) {
	allFields := append(a.fields, args...)
	a.adapter.Debug(msg, allFields...)
}

// Info 记录info级别日志
func (a *LoggerAdapterWithFields) Info(msg string, args ...logger.Field) {
	allFields := append(a.fields, args...)
	a.adapter.Info(msg, allFields...)
}

// Warn 记录warn级别日志
func (a *LoggerAdapterWithFields) Warn(msg string, args ...logger.Field) {
	allFields := append(a.fields, args...)
	a.adapter.Warn(msg, allFields...)
}

// Error 记录error级别日志
func (a *LoggerAdapterWithFields) Error(msg string, args ...logger.Field) {
	allFields := append(a.fields, args...)
	a.adapter.Error(msg, allFields...)
}

// DebugContext 记录debug级别日志，支持context
func (a *LoggerAdapterWithFields) DebugContext(ctx context.Context, msg string, args ...logger.Field) {
	contextFields := FieldsFromContext(ctx)
	allFields := append(a.fields, contextFields...)
	allFields = append(allFields, args...)
	a.adapter.Debug(msg, allFields...)
}

// InfoContext 记录info级别日志，支持context
func (a *LoggerAdapterWithFields) InfoContext(ctx context.Context, msg string, args ...logger.Field) {
	contextFields := FieldsFromContext(ctx)
	allFields := append(a.fields, contextFields...)
	allFields = append(allFields, args...)
	a.adapter.Info(msg, allFields...)
}

// WarnContext 记录warn级别日志，支持context
func (a *LoggerAdapterWithFields) WarnContext(ctx context.Context, msg string, args ...logger.Field) {
	contextFields := FieldsFromContext(ctx)
	allFields := append(a.fields, contextFields...)
	allFields = append(allFields, args...)
	a.adapter.Warn(msg, allFields...)
}

// ErrorContext 记录error级别日志，支持context
func (a *LoggerAdapterWithFields) ErrorContext(ctx context.Context, msg string, args ...logger.Field) {
	contextFields := FieldsFromContext(ctx)
	allFields := append(a.fields, contextFields...)
	allFields = append(allFields, args...)
	a.adapter.Error(msg, allFields...)
}

// With 返回带有更多fields的新适配器
func (a *LoggerAdapterWithFields) With(args ...logger.Field) Logger {
	newFields := append(a.fields, args...)
	return &LoggerAdapterWithFields{
		adapter: a.adapter,
		fields:  newFields,
	}
}

// WithContext 返回带有context的新适配器
func (a *LoggerAdapterWithFields) WithContext(ctx context.Context) Logger {
	contextFields := FieldsFromContext(ctx)
	newFields := append(a.fields, contextFields...)
	return &LoggerAdapterWithFields{
		adapter: a.adapter,
		fields:  newFields,
	}
}
