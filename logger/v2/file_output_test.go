package loggerv2

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/to404hanga/pkg404/logger"
)

// TestOutputModeConstants 测试输出模式常量
func TestOutputModeConstants(t *testing.T) {
	tests := []struct {
		mode     OutputMode
		expected string
	}{
		{OutputModeConsole, "console"},
		{OutputModeFile, "file"},
		{OutputModeBoth, "both"},
	}

	for _, test := range tests {
		if string(test.mode) != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, string(test.mode))
		}
	}
}

// TestOutputConfig 测试输出配置结构
func TestOutputConfig(t *testing.T) {
	config := &OutputConfig{
		Mode:          OutputModeFile,
		FilePath:      "./test.log",
		AutoCreateDir: true,
		Level:         "dev",
	}

	if config.Mode != OutputModeFile {
		t.Errorf("Expected mode %s, got %s", OutputModeFile, config.Mode)
	}
	if config.FilePath != "./test.log" {
		t.Errorf("Expected file path './test.log', got %s", config.FilePath)
	}
	if !config.AutoCreateDir {
		t.Error("Expected AutoCreateDir to be true")
	}
	if config.Level != "dev" {
		t.Errorf("Expected level 'dev', got %s", config.Level)
	}
}

// TestDefaultLoggerConfigWithOutput 测试默认配置包含输出设置
func TestDefaultLoggerConfigWithOutput(t *testing.T) {
	config := DefaultLoggerConfig()

	if config.Output == nil {
		t.Fatal("Expected Output config to be set")
	}
	if config.Output.Mode != OutputModeConsole {
		t.Errorf("Expected default mode %s, got %s", OutputModeConsole, config.Output.Mode)
	}
	if !config.Output.AutoCreateDir {
		t.Error("Expected AutoCreateDir to be true by default")
	}
	if config.Output.Level != "dev" {
		t.Errorf("Expected default level 'dev', got %s", config.Output.Level)
	}
}

// TestDefaultFileLoggerConfig 测试文件输出配置
func TestDefaultFileLoggerConfig(t *testing.T) {
	filePath := "./test/app.log"
	config := DefaultFileLoggerConfig(filePath)

	if config.Output.Mode != OutputModeFile {
		t.Errorf("Expected mode %s, got %s", OutputModeFile, config.Output.Mode)
	}
	if config.Output.FilePath != filePath {
		t.Errorf("Expected file path %s, got %s", filePath, config.Output.FilePath)
	}
}

// TestDefaultBothLoggerConfig 测试同时输出配置
func TestDefaultBothLoggerConfig(t *testing.T) {
	filePath := "./test/app.log"
	config := DefaultBothLoggerConfig(filePath)

	if config.Output.Mode != OutputModeBoth {
		t.Errorf("Expected mode %s, got %s", OutputModeBoth, config.Output.Mode)
	}
	if config.Output.FilePath != filePath {
		t.Errorf("Expected file path %s, got %s", filePath, config.Output.FilePath)
	}
}

// TestEnsureDir 测试目录创建功能
func TestEnsureDir(t *testing.T) {
	// 创建临时测试目录
	testDir := "./test_logs/subdir"
	testFile := filepath.Join(testDir, "test.log")

	// 清理测试目录
	defer func() {
		os.RemoveAll("./test_logs")
	}()

	// 测试目录创建
	err := ensureDir(testFile)
	if err != nil {
		t.Fatalf("Failed to ensure directory: %v", err)
	}

	// 验证目录是否存在
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}

// TestEnsureDirExisting 测试已存在目录的情况
func TestEnsureDirExisting(t *testing.T) {
	// 创建测试目录
	testDir := "./existing_test_logs"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	defer func() {
		os.RemoveAll(testDir)
	}()

	testFile := filepath.Join(testDir, "test.log")

	// 测试已存在目录的情况
	err = ensureDir(testFile)
	if err != nil {
		t.Errorf("ensureDir failed for existing directory: %v", err)
	}
}

// TestNewConsoleLogger 测试控制台 logger 创建
func TestNewConsoleLogger(t *testing.T) {
	logger, err := NewConsoleLogger()
	if err != nil {
		t.Fatalf("Failed to create console logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}
	if logger.config.Output.Mode != OutputModeConsole {
		t.Errorf("Expected console mode, got %s", logger.config.Output.Mode)
	}
}

// TestNewFileLogger 测试文件 logger 创建
func TestNewFileLogger(t *testing.T) {
	testFile := "./test_logs/file_logger_test.log"

	defer func() {
		os.RemoveAll("./test_logs")
	}()

	logger, err := NewFileLogger(testFile)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}
	if logger.config.Output.Mode != OutputModeFile {
		t.Errorf("Expected file mode, got %s", logger.config.Output.Mode)
	}
	if logger.config.Output.FilePath != testFile {
		t.Errorf("Expected file path %s, got %s", testFile, logger.config.Output.FilePath)
	}
}

// TestNewBothLogger 测试同时输出 logger 创建
func TestNewBothLogger(t *testing.T) {
	testFile := "./test_logs/both_logger_test.log"

	l, err := NewBothLogger(testFile)
	if err != nil {
		t.Fatalf("Failed to create both logger: %v", err)
	}

	if l == nil {
		t.Fatal("Logger is nil")
	}
	if l.config.Output.Mode != OutputModeBoth {
		t.Errorf("Expected both mode, got %s", l.config.Output.Mode)
	}

	l.Error("Test error message")
}

// TestNewLoggerWithLevel 测试指定级别的 logger 创建
func TestNewLoggerWithLevel(t *testing.T) {
	testFile := "./test_logs/level_logger_test.log"

	defer func() {
		os.RemoveAll("./test_logs")
	}()

	logger, err := NewLoggerWithLevel(OutputModeFile, testFile, "prod")
	if err != nil {
		t.Fatalf("Failed to create logger with level: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}
	if logger.config.Output.Level != "prod" {
		t.Errorf("Expected level 'prod', got %s", logger.config.Output.Level)
	}
}

// TestFileLoggerFunctionality 测试文件 logger 的实际功能
func TestFileLoggerFunctionality(t *testing.T) {
	testFile := "./test_logs/functionality_test.log"

	defer func() {
		os.RemoveAll("./test_logs")
	}()

	l, err := NewFileLogger(testFile)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	// 测试基本日志记录
	l.Info("Test info message")
	l.Debug("Test debug message")
	l.Warn("Test warn message")
	l.Error("Test error message")

	// 测试带字段的日志记录
	l.WithFields(
		logger.String("key1", "value1"),
		logger.Int("key2", 42),
	).Info("Test message with fields")

	// 测试上下文日志记录
	ctx := context.WithValue(context.Background(), "user_id", "12345")
	l.WithContext(ctx).InfoContext(ctx, "Test context message")

	// 验证文件是否存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

// TestNewZapCtxLoggerWithConfig 测试使用配置创建 logger
func TestNewZapCtxLoggerWithConfig(t *testing.T) {
	config := &LoggerConfig{
		EnableContextExtraction: true,
		ContextExtractor: &DefaultContextExtractor{
			Keys: []any{"user_id", "request_id"},
		},
		ExtractAllContextValues: false,
		Output: &OutputConfig{
			Mode:          OutputModeFile,
			FilePath:      "./test_logs/config_test.log",
			AutoCreateDir: true,
			Level:         "dev",
		},
	}

	defer func() {
		os.RemoveAll("./test_logs")
	}()

	logger, err := NewZapCtxLoggerWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to create logger with config: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}

	// 测试配置是否正确设置
	if logger.config.Output.Mode != OutputModeFile {
		t.Errorf("Expected file mode, got %s", logger.config.Output.Mode)
	}
	if !logger.config.EnableContextExtraction {
		t.Error("Expected context extraction to be enabled")
	}
}

// TestBuildZapConfig 测试 zap 配置构建
func TestBuildZapConfig(t *testing.T) {
	tests := []struct {
		name   string
		output *OutputConfig
	}{
		{
			name: "Console output",
			output: &OutputConfig{
				Mode:          OutputModeConsole,
				AutoCreateDir: true,
				Level:         "dev",
			},
		},
		{
			name: "File output",
			output: &OutputConfig{
				Mode:          OutputModeFile,
				FilePath:      "./test_logs/build_config_test.log",
				AutoCreateDir: true,
				Level:         "prod",
			},
		},
		{
			name: "Both output",
			output: &OutputConfig{
				Mode:          OutputModeBoth,
				FilePath:      "./test_logs/build_config_both_test.log",
				AutoCreateDir: true,
				Level:         "test",
			},
		},
	}

	defer func() {
		os.RemoveAll("./test_logs")
	}()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := buildZapConfig(test.output)
			if err != nil {
				t.Fatalf("Failed to build zap config: %v", err)
			}

			// 验证输出路径设置
			switch test.output.Mode {
			case OutputModeConsole:
				if len(cfg.OutputPaths) != 1 || cfg.OutputPaths[0] != "stdout" {
					t.Errorf("Expected stdout output, got %v", cfg.OutputPaths)
				}
			case OutputModeFile:
				if len(cfg.OutputPaths) != 1 || cfg.OutputPaths[0] != test.output.FilePath {
					t.Errorf("Expected file output %s, got %v", test.output.FilePath, cfg.OutputPaths)
				}
			case OutputModeBoth:
				if len(cfg.OutputPaths) != 2 {
					t.Errorf("Expected 2 outputs, got %d", len(cfg.OutputPaths))
				}
			}
		})
	}
}

// TestAutoDirectoryCreation 测试自动目录创建功能
func TestAutoDirectoryCreation(t *testing.T) {
	// 使用深层嵌套的目录路径
	testFile := "./test_logs/deep/nested/directory/auto_create_test.log"

	defer func() {
		os.RemoveAll("./test_logs")
	}()

	config := &OutputConfig{
		Mode:          OutputModeFile,
		FilePath:      testFile,
		AutoCreateDir: true,
		Level:         "dev",
	}

	logger, err := NewZapCtxLoggerWithConfig(&LoggerConfig{
		Output: config,
	})
	if err != nil {
		t.Fatalf("Failed to create logger with auto directory creation: %v", err)
	}

	// 测试日志记录
	logger.Info("Test message for auto directory creation")

	// 验证目录和文件是否存在
	if _, err := os.Stat(filepath.Dir(testFile)); os.IsNotExist(err) {
		t.Error("Auto-created directory does not exist")
	}
}
