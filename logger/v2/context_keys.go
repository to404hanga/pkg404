package loggerv2

// 定义常用的context key类型，避免字符串冲突
type contextKey string

// 常用的context key定义
const (
	// 用户相关
	UserIDKey   contextKey = "user_id"
	UserNameKey contextKey = "user_name"
	TenantIDKey contextKey = "tenant_id"
	ClientIDKey contextKey = "client_id"

	// 请求相关
	RequestIDKey     contextKey = "request_id"
	CorrelationIDKey contextKey = "correlation_id"
	SessionIDKey     contextKey = "session_id"

	// 链路追踪相关
	TraceIDKey contextKey = "trace_id"
	SpanIDKey  contextKey = "span_id"

	// 网络相关
	IPKey        contextKey = "ip"
	UserAgentKey contextKey = "user_agent"

	// 业务相关
	OperationKey contextKey = "operation"
	ModuleKey    contextKey = "module"
	ServiceKey   contextKey = "service"
)

// String 返回context key的字符串表示
func (c contextKey) String() string {
	return string(c)
}

// CommonContextKeys 返回所有常用的context key
func CommonContextKeys() []any {
	return []any{
		UserIDKey, UserNameKey, TenantIDKey, ClientIDKey,
		RequestIDKey, CorrelationIDKey, SessionIDKey,
		TraceIDKey, SpanIDKey,
		IPKey, UserAgentKey,
		OperationKey, ModuleKey, ServiceKey,
	}
}

// NewCommonContextExtractor 创建包含所有常用key的context提取器
func NewCommonContextExtractor() *DefaultContextExtractor {
	return &DefaultContextExtractor{
		Keys: CommonContextKeys(),
	}
}

// NewCustomContextExtractor 创建自定义key的context提取器
func NewCustomContextExtractor(keys ...any) *DefaultContextExtractor {
	return &DefaultContextExtractor{
		Keys: keys,
	}
}
