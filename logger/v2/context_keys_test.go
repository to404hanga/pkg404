package loggerv2

import (
	"context"
	"testing"
)

// TestContextKeyString 测试 contextKey 的 String 方法
func TestContextKeyString(t *testing.T) {
	key := contextKey("test_key")
	expected := "test_key"
	
	if key.String() != expected {
		t.Errorf("expected %s, got %s", expected, key.String())
	}
}

// TestContextKeyConstants 测试所有预定义的 context key 常量
func TestContextKeyConstants(t *testing.T) {
	tests := []struct {
		key      contextKey
		expected string
	}{
		{UserIDKey, "user_id"},
		{UserNameKey, "user_name"},
		{TenantIDKey, "tenant_id"},
		{ClientIDKey, "client_id"},
		{RequestIDKey, "request_id"},
		{CorrelationIDKey, "correlation_id"},
		{SessionIDKey, "session_id"},
		{TraceIDKey, "trace_id"},
		{SpanIDKey, "span_id"},
		{IPKey, "ip"},
		{UserAgentKey, "user_agent"},
		{OperationKey, "operation"},
		{ModuleKey, "module"},
		{ServiceKey, "service"},
	}

	for _, test := range tests {
		if string(test.key) != test.expected {
			t.Errorf("expected %s, got %s", test.expected, string(test.key))
		}
		
		if test.key.String() != test.expected {
			t.Errorf("String() method: expected %s, got %s", test.expected, test.key.String())
		}
	}
}

// TestCommonContextKeys 测试 CommonContextKeys 函数
func TestCommonContextKeys(t *testing.T) {
	keys := CommonContextKeys()
	
	if len(keys) == 0 {
		t.Error("CommonContextKeys should return non-empty slice")
	}
	
	// 验证返回的 keys 包含所有预定义的常量
	expectedKeys := []contextKey{
		UserIDKey, UserNameKey, TenantIDKey, ClientIDKey,
		RequestIDKey, CorrelationIDKey, SessionIDKey,
		TraceIDKey, SpanIDKey,
		IPKey, UserAgentKey,
		OperationKey, ModuleKey, ServiceKey,
	}
	
	if len(keys) != len(expectedKeys) {
		t.Errorf("expected %d keys, got %d", len(expectedKeys), len(keys))
	}
	
	// 将返回的 keys 转换为 map 以便查找
	keyMap := make(map[contextKey]bool)
	for _, key := range keys {
		if contextKey, ok := key.(contextKey); ok {
			keyMap[contextKey] = true
		}
	}
	
	// 验证所有预期的 key 都存在
	for _, expectedKey := range expectedKeys {
		if !keyMap[expectedKey] {
			t.Errorf("expected key %s not found in CommonContextKeys", expectedKey)
		}
	}
}

// TestNewCommonContextExtractor 测试 NewCommonContextExtractor 函数
func TestNewCommonContextExtractor(t *testing.T) {
	extractor := NewCommonContextExtractor()
	
	if extractor == nil {
		t.Fatal("NewCommonContextExtractor should not return nil")
	}
	
	if len(extractor.Keys) == 0 {
		t.Error("common context extractor should have keys")
	}
	
	// 验证返回的提取器包含所有常用 key
	commonKeys := CommonContextKeys()
	if len(extractor.Keys) != len(commonKeys) {
		t.Errorf("expected %d keys, got %d", len(commonKeys), len(extractor.Keys))
	}
	
	// 验证 keys 的内容一致
	for i, key := range commonKeys {
		if extractor.Keys[i] != key {
			t.Errorf("key %d mismatch: expected %v, got %v", i, key, extractor.Keys[i])
		}
	}
}

// TestNewCustomContextExtractor 测试 NewCustomContextExtractor 函数
func TestNewCustomContextExtractor(t *testing.T) {
	customKeys := []any{"custom1", "custom2", UserIDKey}
	extractor := NewCustomContextExtractor(customKeys...)
	
	if extractor == nil {
		t.Fatal("NewCustomContextExtractor should not return nil")
	}
	
	if len(extractor.Keys) != len(customKeys) {
		t.Errorf("expected %d keys, got %d", len(customKeys), len(extractor.Keys))
	}
	
	// 验证 keys 的内容一致
	for i, key := range customKeys {
		if extractor.Keys[i] != key {
			t.Errorf("key %d mismatch: expected %v, got %v", i, key, extractor.Keys[i])
		}
	}
}

// TestNewCustomContextExtractorEmpty 测试空参数的 NewCustomContextExtractor
func TestNewCustomContextExtractorEmpty(t *testing.T) {
	extractor := NewCustomContextExtractor()
	
	if extractor == nil {
		t.Fatal("NewCustomContextExtractor should not return nil")
	}
	
	if len(extractor.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(extractor.Keys))
	}
}

// TestContextKeyUsageInContext 测试在实际 context 中使用 contextKey
func TestContextKeyUsageInContext(t *testing.T) {
	ctx := context.Background()
	
	// 使用 contextKey 设置值
	ctx = context.WithValue(ctx, UserIDKey, "user123")
	ctx = context.WithValue(ctx, RequestIDKey, "req456")
	ctx = context.WithValue(ctx, TraceIDKey, "trace789")
	
	// 验证可以正确获取值
	if userID := ctx.Value(UserIDKey); userID != "user123" {
		t.Errorf("expected user123, got %v", userID)
	}
	
	if requestID := ctx.Value(RequestIDKey); requestID != "req456" {
		t.Errorf("expected req456, got %v", requestID)
	}
	
	if traceID := ctx.Value(TraceIDKey); traceID != "trace789" {
		t.Errorf("expected trace789, got %v", traceID)
	}
	
	// 验证不存在的 key 返回 nil
	if spanID := ctx.Value(SpanIDKey); spanID != nil {
		t.Errorf("expected nil for SpanIDKey, got %v", spanID)
	}
}

// TestContextKeyTypesSafety 测试 contextKey 类型安全性
func TestContextKeyTypesSafety(t *testing.T) {
	ctx := context.Background()
	
	// 使用字符串 key 设置值
	ctx = context.WithValue(ctx, "user_id", "string_key_value")
	
	// 使用 contextKey 设置值
	ctx = context.WithValue(ctx, UserIDKey, "context_key_value")
	
	// 验证两者是不同的 key
	stringValue := ctx.Value("user_id")
	contextKeyValue := ctx.Value(UserIDKey)
	
	if stringValue != "string_key_value" {
		t.Errorf("expected string_key_value, got %v", stringValue)
	}
	
	if contextKeyValue != "context_key_value" {
		t.Errorf("expected context_key_value, got %v", contextKeyValue)
	}
	
	// 验证它们确实是不同的值
	if stringValue == contextKeyValue {
		t.Error("string key and contextKey should be different")
	}
}

// TestContextKeyExtractorIntegration 测试 contextKey 与提取器的集成
func TestContextKeyExtractorIntegration(t *testing.T) {
	// 创建包含 contextKey 的提取器
	extractor := NewCommonContextExtractor()
	
	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, "user123")
	ctx = context.WithValue(ctx, RequestIDKey, "req456")
	ctx = context.WithValue(ctx, "string_key", "string_value") // 字符串 key，不会被提取
	
	fields := extractor.ExtractFields(ctx)
	
	// 应该提取到 2 个字段
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
	
	// 验证提取的字段
	fieldMap := make(map[string]interface{})
	for _, field := range fields {
		fieldMap[field.Key] = field.Val
	}
	
	if fieldMap["user_id"] != "user123" {
		t.Errorf("expected user123, got %v", fieldMap["user_id"])
	}
	
	if fieldMap["request_id"] != "req456" {
		t.Errorf("expected req456, got %v", fieldMap["request_id"])
	}
	
	// 验证字符串 key 不被提取
	if _, exists := fieldMap["string_key"]; exists {
		t.Error("string_key should not be extracted")
	}
}

// TestContextKeyUniqueness 测试 contextKey 的唯一性
func TestContextKeyUniqueness(t *testing.T) {
	// 验证所有预定义的 contextKey 都是唯一的
	keys := []contextKey{
		UserIDKey, UserNameKey, TenantIDKey, ClientIDKey,
		RequestIDKey, CorrelationIDKey, SessionIDKey,
		TraceIDKey, SpanIDKey,
		IPKey, UserAgentKey,
		OperationKey, ModuleKey, ServiceKey,
	}
	
	keySet := make(map[string]bool)
	for _, key := range keys {
		keyStr := string(key)
		if keySet[keyStr] {
			t.Errorf("duplicate key found: %s", keyStr)
		}
		keySet[keyStr] = true
	}
	
	// 验证总数正确
	if len(keySet) != len(keys) {
		t.Errorf("expected %d unique keys, got %d", len(keys), len(keySet))
	}
}