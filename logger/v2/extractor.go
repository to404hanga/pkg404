package loggerv2

import (
	"context"

	"github.com/to404hanga/pkg404/logger"
)

// ContextExtractor 定义如何从 context 中提取值的接口
type ContextExtractor interface {
	// ExtractFields 从 context 中提取字段
	ExtractFields(ctx context.Context) []logger.Field
}

// DefaultContextExtractor 默认的 context 提取器
type DefaultContextExtractor struct {
	// 需要提取的 context key 列表
	Keys []any
}

// ExtractFields 实现 ContextExtractor 接口，从 context 中提取指定 key 的值
func (e *DefaultContextExtractor) ExtractFields(ctx context.Context) []logger.Field {
	if ctx == nil {
		return nil
	}

	var fields []logger.Field
	for _, key := range e.Keys {
		if value := ctx.Value(key); value != nil {
			// 将 key 转换为字符串作为字段名
			keyStr := ""
			switch k := key.(type) {
			case string:
				keyStr = k
			case contextKey:
				keyStr = string(k)
			case int:
				keyStr = "ctx_key_int"
			default:
				keyStr = "ctx_value"
			}
			fields = append(fields, logger.Field{Key: keyStr, Val: value})
		}
	}
	return fields
}

// AllContextExtractor 提取 context 中所有可访问的值（通过反射）
type AllContextExtractor struct{}

// ExtractFields 提取 context 中所有可访问的值
func (e *AllContextExtractor) ExtractFields(ctx context.Context) []logger.Field {
	if ctx == nil {
		return nil
	}

	var fields []logger.Field

	// 遍历 context 的值（这里使用一种安全的方式）
	// 由于 Go 的 context 包没有直接暴露所有值的方法，这里使用一些常见的 key
	commonKeys := []any{
		"user_id", "request_id", "trace_id", "span_id", "session_id",
		"correlation_id", "tenant_id", "client_id", "ip", "user_agent",
	}

	for _, key := range commonKeys {
		if value := ctx.Value(key); value != nil {
			keyStr := ""
			if s, ok := key.(string); ok {
				keyStr = s
			} else {
				keyStr = "ctx_value"
			}
			fields = append(fields, logger.Field{Key: keyStr, Val: value})
		}
	}

	return fields
}
