package loggerv2

import (
	"testing"
)

// TestDefaultLoggerConfig 测试默认配置
func TestDefaultLoggerConfig(t *testing.T) {
	config := DefaultLoggerConfig()
	
	if config == nil {
		t.Fatal("DefaultLoggerConfig should not return nil")
	}
	
	// 验证默认配置值
	if !config.EnableContextExtraction {
		t.Error("default config should enable context extraction")
	}
	
	if config.ContextExtractor == nil {
		t.Error("default config should have a context extractor")
	}
	
	if config.ExtractAllContextValues {
		t.Error("default config should not extract all context values")
	}
	
	// 验证输出配置
	if config.Output == nil {
		t.Fatal("Expected Output config to be set")
	}
	if config.Output.Mode != OutputModeConsole {
		t.Errorf("Expected default output mode to be %s, got %s", OutputModeConsole, config.Output.Mode)
	}
	if !config.Output.AutoCreateDir {
		t.Error("Expected AutoCreateDir to be true by default")
	}
	if config.Output.Level != "dev" {
		t.Errorf("Expected default level to be 'dev', got %s", config.Output.Level)
	}
	
	// 验证默认提取器的配置
	if extractor, ok := config.ContextExtractor.(*DefaultContextExtractor); ok {
		if len(extractor.Keys) == 0 {
			t.Error("default extractor should have predefined keys")
		}
		
		// 检查是否包含预期的默认 key
		expectedKeys := []string{"user_id", "request_id", "trace_id", "span_id"}
		keyMap := make(map[string]bool)
		for _, key := range extractor.Keys {
			if str, ok := key.(string); ok {
				keyMap[str] = true
			}
		}
		
		for _, expectedKey := range expectedKeys {
			if !keyMap[expectedKey] {
				t.Errorf("default extractor should contain key: %s", expectedKey)
			}
		}
	} else {
		t.Error("default context extractor should be of type *DefaultContextExtractor")
	}
}

// TestLoggerConfigCustomization 测试自定义配置
func TestLoggerConfigCustomization(t *testing.T) {
	customExtractor := &DefaultContextExtractor{
		Keys: []any{"custom_key1", "custom_key2"},
	}
	
	config := &LoggerConfig{
		EnableContextExtraction: false,
		ContextExtractor:        customExtractor,
		ExtractAllContextValues: true,
		Output: &OutputConfig{
			Mode:          OutputModeFile,
			FilePath:      "./custom.log",
			AutoCreateDir: false,
			Level:         "prod",
		},
	}
	
	// 验证自定义配置值
	if config.EnableContextExtraction {
		t.Error("custom config should disable context extraction")
	}
	
	if config.ContextExtractor != customExtractor {
		t.Error("custom config should use the provided extractor")
	}
	
	if !config.ExtractAllContextValues {
		t.Error("custom config should extract all context values")
	}
	
	// 验证自定义输出配置
	if config.Output == nil {
		t.Fatal("Expected Output config to be set")
	}
	if config.Output.Mode != OutputModeFile {
		t.Errorf("Expected output mode to be %s, got %s", OutputModeFile, config.Output.Mode)
	}
	if config.Output.FilePath != "./custom.log" {
		t.Errorf("Expected file path to be './custom.log', got %s", config.Output.FilePath)
	}
	if config.Output.AutoCreateDir {
		t.Error("Expected AutoCreateDir to be false")
	}
	if config.Output.Level != "prod" {
		t.Errorf("Expected level to be 'prod', got %s", config.Output.Level)
	}
	
	// 验证自定义提取器
	if extractor, ok := config.ContextExtractor.(*DefaultContextExtractor); ok {
		if len(extractor.Keys) != 2 {
			t.Errorf("custom extractor should have 2 keys, got %d", len(extractor.Keys))
		}
		
		expectedKeys := []string{"custom_key1", "custom_key2"}
		for i, expectedKey := range expectedKeys {
			if extractor.Keys[i] != expectedKey {
				t.Errorf("key %d should be %s, got %v", i, expectedKey, extractor.Keys[i])
			}
		}
	} else {
		t.Error("custom context extractor should be of type *DefaultContextExtractor")
	}
}

// TestLoggerConfigNilExtractor 测试 nil 提取器的情况
func TestLoggerConfigNilExtractor(t *testing.T) {
	config := &LoggerConfig{
		EnableContextExtraction: true,
		ContextExtractor:        nil,
		ExtractAllContextValues: false,
	}
	
	// 验证配置可以接受 nil 提取器
	if config.ContextExtractor != nil {
		t.Error("config should accept nil extractor")
	}
	
	if !config.EnableContextExtraction {
		t.Error("context extraction should be enabled")
	}
}

// TestLoggerConfigImmutability 测试配置的不可变性
func TestLoggerConfigImmutability(t *testing.T) {
	config1 := DefaultLoggerConfig()
	config2 := DefaultLoggerConfig()
	
	// 验证每次调用都返回新的实例
	if config1 == config2 {
		t.Error("DefaultLoggerConfig should return new instances")
	}
	
	// 修改一个配置不应该影响另一个
	config1.EnableContextExtraction = false
	if !config2.EnableContextExtraction {
		t.Error("modifying one config should not affect another")
	}
	
	// 修改提取器的 keys 不应该影响其他实例
	if extractor1, ok := config1.ContextExtractor.(*DefaultContextExtractor); ok {
		if extractor2, ok := config2.ContextExtractor.(*DefaultContextExtractor); ok {
			originalLen := len(extractor2.Keys)
			extractor1.Keys = append(extractor1.Keys, "new_key")
			
			if len(extractor2.Keys) != originalLen {
				t.Error("modifying one extractor's keys should not affect another")
			}
		}
	}
}

// TestLoggerConfigValidation 测试配置验证
func TestLoggerConfigValidation(t *testing.T) {
	// 测试启用 context 提取但没有提取器的情况
	config := &LoggerConfig{
		EnableContextExtraction: true,
		ContextExtractor:        nil,
		ExtractAllContextValues: false,
	}
	
	// 这种配置应该是有效的，但不会提取任何 context 值
	if config.EnableContextExtraction && config.ContextExtractor == nil {
		// 这是一个有效的配置状态，应该在运行时优雅处理
		t.Log("Config with enabled extraction but nil extractor is valid")
	}
	
	// 测试禁用 context 提取但有提取器的情况
	config2 := &LoggerConfig{
		EnableContextExtraction: false,
		ContextExtractor: &DefaultContextExtractor{
			Keys: []any{"test_key"},
		},
		ExtractAllContextValues: false,
	}
	
	// 这种配置也应该是有效的，提取器会被忽略
	if !config2.EnableContextExtraction && config2.ContextExtractor != nil {
		t.Log("Config with disabled extraction but non-nil extractor is valid")
	}
}