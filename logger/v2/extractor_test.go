package loggerv2

import (
	"context"
	"testing"
)

// TestDefaultContextExtractor 测试默认 context 提取器
func TestDefaultContextExtractor(t *testing.T) {
	extractor := &DefaultContextExtractor{
		Keys: []any{"user_id", "request_id", "trace_id"},
	}

	// 测试 nil context
	fields := extractor.ExtractFields(nil)
	if fields != nil {
		t.Error("ExtractFields should return nil for nil context")
	}

	// 测试空 context
	ctx := context.Background()
	fields = extractor.ExtractFields(ctx)
	if len(fields) != 0 {
		t.Errorf("expected 0 fields for empty context, got %d", len(fields))
	}

	// 测试包含值的 context
	ctx = context.WithValue(ctx, "user_id", "12345")
	ctx = context.WithValue(ctx, "request_id", "req-67890")
	ctx = context.WithValue(ctx, "unknown_key", "unknown_value")

	fields = extractor.ExtractFields(ctx)
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}

	// 验证提取的字段
	fieldMap := make(map[string]interface{})
	for _, field := range fields {
		fieldMap[field.Key] = field.Val
	}

	if fieldMap["user_id"] != "12345" {
		t.Errorf("expected user_id to be '12345', got %v", fieldMap["user_id"])
	}

	if fieldMap["request_id"] != "req-67890" {
		t.Errorf("expected request_id to be 'req-67890', got %v", fieldMap["request_id"])
	}

	// 验证未在 Keys 中的值不被提取
	if _, exists := fieldMap["unknown_key"]; exists {
		t.Error("unknown_key should not be extracted")
	}
}

// TestDefaultContextExtractorWithContextKey 测试使用 contextKey 类型的提取器
func TestDefaultContextExtractorWithContextKey(t *testing.T) {
	extractor := &DefaultContextExtractor{
		Keys: []any{UserIDKey, RequestIDKey, TraceIDKey},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIDKey, "user123")
	ctx = context.WithValue(ctx, RequestIDKey, "req456")
	ctx = context.WithValue(ctx, SpanIDKey, "span789") // 不在 Keys 中

	fields := extractor.ExtractFields(ctx)
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}

	// 验证提取的字段
	fieldMap := make(map[string]interface{})
	for _, field := range fields {
		fieldMap[field.Key] = field.Val
	}

	if fieldMap["user_id"] != "user123" {
		t.Errorf("expected user_id to be 'user123', got %v", fieldMap["user_id"])
	}

	if fieldMap["request_id"] != "req456" {
		t.Errorf("expected request_id to be 'req456', got %v", fieldMap["request_id"])
	}

	// 验证 SpanIDKey 不被提取
	if _, exists := fieldMap["span_id"]; exists {
		t.Error("span_id should not be extracted")
	}
}

// TestDefaultContextExtractorWithMixedKeys 测试混合类型的 key
func TestDefaultContextExtractorWithMixedKeys(t *testing.T) {
	extractor := &DefaultContextExtractor{
		Keys: []any{"string_key", UserIDKey, 42, struct{}{}},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "string_key", "string_value")
	ctx = context.WithValue(ctx, UserIDKey, "user_value")
	ctx = context.WithValue(ctx, 42, "int_value")
	ctx = context.WithValue(ctx, struct{}{}, "struct_value")

	fields := extractor.ExtractFields(ctx)
	if len(fields) != 4 {
		t.Errorf("expected 4 fields, got %d", len(fields))
	}

	// 验证不同类型的 key 都能正确处理
	fieldMap := make(map[string]interface{})
	for _, field := range fields {
		fieldMap[field.Key] = field.Val
	}

	if fieldMap["string_key"] != "string_value" {
		t.Errorf("expected string_key to be 'string_value', got %v", fieldMap["string_key"])
	}

	if fieldMap["user_id"] != "user_value" {
		t.Errorf("expected user_id to be 'user_value', got %v", fieldMap["user_id"])
	}

	if fieldMap["ctx_key_int"] != "int_value" {
		t.Errorf("expected ctx_key_int to be 'int_value', got %v", fieldMap["ctx_key_int"])
	}

	if fieldMap["ctx_value"] != "struct_value" {
		t.Errorf("expected ctx_value to be 'struct_value', got %v", fieldMap["ctx_value"])
	}
}

// TestAllContextExtractor 测试全量 context 提取器
func TestAllContextExtractor(t *testing.T) {
	extractor := &AllContextExtractor{}

	// 测试 nil context
	fields := extractor.ExtractFields(nil)
	if fields != nil {
		t.Error("ExtractFields should return nil for nil context")
	}

	// 测试空 context
	ctx := context.Background()
	fields = extractor.ExtractFields(ctx)
	if len(fields) != 0 {
		t.Errorf("expected 0 fields for empty context, got %d", len(fields))
	}

	// 测试包含常见 key 的 context
	ctx = context.WithValue(ctx, "user_id", "12345")
	ctx = context.WithValue(ctx, "request_id", "req-67890")
	ctx = context.WithValue(ctx, "trace_id", "trace-111")
	ctx = context.WithValue(ctx, "unknown_key", "unknown_value") // 不在常见 key 列表中

	fields = extractor.ExtractFields(ctx)

	// AllContextExtractor 只提取预定义的常见 key
	if len(fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fields))
	}

	// 验证提取的字段
	fieldMap := make(map[string]interface{})
	for _, field := range fields {
		fieldMap[field.Key] = field.Val
	}

	if fieldMap["user_id"] != "12345" {
		t.Errorf("expected user_id to be '12345', got %v", fieldMap["user_id"])
	}

	if fieldMap["request_id"] != "req-67890" {
		t.Errorf("expected request_id to be 'req-67890', got %v", fieldMap["request_id"])
	}

	if fieldMap["trace_id"] != "trace-111" {
		t.Errorf("expected trace_id to be 'trace-111', got %v", fieldMap["trace_id"])
	}

	// 验证不在常见 key 列表中的值不被提取
	if _, exists := fieldMap["unknown_key"]; exists {
		t.Error("unknown_key should not be extracted by AllContextExtractor")
	}
}

// TestContextExtractorInterface 测试 ContextExtractor 接口
func TestContextExtractorInterface(t *testing.T) {
	// 验证 DefaultContextExtractor 实现了接口
	var extractor ContextExtractor = &DefaultContextExtractor{
		Keys: []any{"test_key"},
	}

	ctx := context.WithValue(context.Background(), "test_key", "test_value")
	fields := extractor.ExtractFields(ctx)

	if len(fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(fields))
	}

	if fields[0].Key != "test_key" || fields[0].Val != "test_value" {
		t.Errorf("expected field {test_key: test_value}, got {%s: %v}", fields[0].Key, fields[0].Val)
	}

	// 验证 AllContextExtractor 实现了接口
	var allExtractor ContextExtractor = &AllContextExtractor{}
	fields = allExtractor.ExtractFields(ctx)

	// AllContextExtractor 不会提取 "test_key"，因为它不在预定义列表中
	if len(fields) != 0 {
		t.Errorf("expected 0 fields from AllContextExtractor, got %d", len(fields))
	}
}

// TestExtractorWithNilValues 测试提取器处理 nil 值的情况
func TestExtractorWithNilValues(t *testing.T) {
	extractor := &DefaultContextExtractor{
		Keys: []any{"nil_key", "valid_key"},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "nil_key", nil)
	ctx = context.WithValue(ctx, "valid_key", "valid_value")

	fields := extractor.ExtractFields(ctx)

	// nil 值不应该被提取
	if len(fields) != 1 {
		t.Errorf("expected 1 field (nil values should be skipped), got %d", len(fields))
	}

	if fields[0].Key != "valid_key" || fields[0].Val != "valid_value" {
		t.Errorf("expected field {valid_key: valid_value}, got {%s: %v}", fields[0].Key, fields[0].Val)
	}
}

// TestExtractorWithEmptyKeys 测试空 Keys 列表的提取器
func TestExtractorWithEmptyKeys(t *testing.T) {
	extractor := &DefaultContextExtractor{
		Keys: []any{},
	}

	ctx := context.WithValue(context.Background(), "some_key", "some_value")
	fields := extractor.ExtractFields(ctx)

	if len(fields) != 0 {
		t.Errorf("expected 0 fields for empty Keys list, got %d", len(fields))
	}
}

// TestExtractorFieldTypes 测试提取器处理不同类型值的情况
func TestExtractorFieldTypes(t *testing.T) {
	extractor := &DefaultContextExtractor{
		Keys: []any{"string", "int", "bool", "slice", "map"},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "string", "hello")
	ctx = context.WithValue(ctx, "int", 42)
	ctx = context.WithValue(ctx, "bool", true)
	ctx = context.WithValue(ctx, "slice", []string{"a", "b", "c"})
	ctx = context.WithValue(ctx, "map", map[string]int{"x": 1, "y": 2})

	fields := extractor.ExtractFields(ctx)

	if len(fields) != 5 {
		t.Errorf("expected 5 fields, got %d", len(fields))
	}

	// 验证所有类型的值都能正确提取
	fieldMap := make(map[string]interface{})
	for _, field := range fields {
		fieldMap[field.Key] = field.Val
	}

	if fieldMap["string"] != "hello" {
		t.Errorf("expected string field to be 'hello', got %v", fieldMap["string"])
	}

	if fieldMap["int"] != 42 {
		t.Errorf("expected int field to be 42, got %v", fieldMap["int"])
	}

	if fieldMap["bool"] != true {
		t.Errorf("expected bool field to be true, got %v", fieldMap["bool"])
	}

	// 验证复杂类型也能正确提取
	if slice, ok := fieldMap["slice"].([]string); !ok || len(slice) != 3 {
		t.Errorf("expected slice field to be []string with 3 elements, got %v", fieldMap["slice"])
	}

	if mapVal, ok := fieldMap["map"].(map[string]int); !ok || len(mapVal) != 2 {
		t.Errorf("expected map field to be map[string]int with 2 elements, got %v", fieldMap["map"])
	}
}
