package loggerv2

// OutputMode 定义日志输出模式
type OutputMode string

const (
	// OutputModeConsole 仅输出到控制台
	OutputModeConsole OutputMode = "console"
	// OutputModeFile 仅输出到文件
	OutputModeFile OutputMode = "file"
	// OutputModeBoth 同时输出到控制台和文件
	OutputModeBoth OutputMode = "both"
)

// OutputConfig 输出配置
type OutputConfig struct {
	// 输出模式
	Mode OutputMode
	// 文件路径（当 Mode 为 file 或 both 时使用）
	FilePath string
	// 是否自动创建目录
	AutoCreateDir bool
	// 日志级别（dev/test/prod）
	Level string
}

// LoggerV2Config v2 版本 logger 的配置
type LoggerConfig struct {
	// 是否启用 context 值提取
	EnableContextExtraction bool
	// context 提取器
	ContextExtractor ContextExtractor
	// 是否自动提取所有 context 值
	ExtractAllContextValues bool
	// 输出配置
	Output *OutputConfig
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
		Output: &OutputConfig{
			Mode:          OutputModeConsole,
			FilePath:      "",
			AutoCreateDir: true,
			Level:         "dev",
		},
	}
}

// DefaultFileLoggerConfig 返回默认的文件输出配置
func DefaultFileLoggerConfig(filePath string) *LoggerConfig {
	config := DefaultLoggerConfig()
	config.Output = &OutputConfig{
		Mode:          OutputModeFile,
		FilePath:      filePath,
		AutoCreateDir: true,
		Level:         "dev",
	}
	return config
}

// DefaultBothLoggerConfig 返回同时输出到控制台和文件的配置
func DefaultBothLoggerConfig(filePath string) *LoggerConfig {
	config := DefaultLoggerConfig()
	config.Output = &OutputConfig{
		Mode:          OutputModeBoth,
		FilePath:      filePath,
		AutoCreateDir: true,
		Level:         "dev",
	}
	return config
}
