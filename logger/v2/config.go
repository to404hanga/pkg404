package loggerv2

// LoggerV2Config v2 版本 logger 的配置
type LoggerConfig struct {
	// 是否启用 context 值提取
	EnableContextExtraction bool
	// context 提取器
	ContextExtractor ContextExtractor
	// 是否自动提取所有 context 值
	ExtractAllContextValues bool
}

// DefaultLoggerConfig 返回默认的 v2 配置
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		EnableContextExtraction: true,
		ContextExtractor: &DefaultContextExtractor{
			Keys: []any{
				"user_id", "request_id", "trace_id", "span_id",
			},
		},
		ExtractAllContextValues: false,
	}
}
