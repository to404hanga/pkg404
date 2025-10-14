package loggerv2

import (
	"context"

	"github.com/to404hanga/pkg404/logger"
	"go.uber.org/zap"
)

// ZapCtxLogger 是 ZapLogger 的 v2 版本，支持 context 值提取
type ZapCtxLogger struct {
	logger     *zap.Logger
	withFields []logger.Field
	config     *LoggerConfig
	ctx        context.Context
}

var _ Logger = (*ZapCtxLogger)(nil)

// NewZapCtxLogger 创建一个新的 ZapCtxLogger 实例
func NewZapCtxLogger(logger *zap.Logger, config ...*LoggerConfig) *ZapCtxLogger {
	cfg := DefaultLoggerConfig()
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	return &ZapCtxLogger{
		logger: logger,
		config: cfg,
	}
}

// NewZapCtxLoggerWithConfig 根据配置创建一个新的 ZapCtxLogger 实例
// 如果配置中包含输出设置，会自动创建相应的 zap.Logger
func NewZapCtxLoggerWithConfig(config *LoggerConfig) (*ZapCtxLogger, error) {
	if config == nil {
		config = DefaultLoggerConfig()
	}
	
	// 如果配置中有输出设置，创建相应的 zap.Logger
	if config.Output != nil {
		zapLogger, err := createZapLogger(config.Output)
		if err != nil {
			return nil, err
		}
		
		return &ZapCtxLogger{
			logger: zapLogger,
			config: config,
		}, nil
	}
	
	// 如果没有输出配置，使用默认的开发配置
	zapConfig := zap.NewDevelopmentConfig()
	zapLogger, err := zapConfig.Build(
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(2),
	)
	if err != nil {
		return nil, err
	}
	
	return &ZapCtxLogger{
		logger: zapLogger,
		config: config,
	}, nil
}

// NewConsoleLogger 创建一个输出到控制台的 logger
func NewConsoleLogger() (*ZapCtxLogger, error) {
	config := DefaultLoggerConfig()
	config.Output.Mode = OutputModeConsole
	return NewZapCtxLoggerWithConfig(config)
}

// NewFileLogger 创建一个输出到文件的 logger
func NewFileLogger(filePath string) (*ZapCtxLogger, error) {
	config := DefaultFileLoggerConfig(filePath)
	return NewZapCtxLoggerWithConfig(config)
}

// NewBothLogger 创建一个同时输出到控制台和文件的 logger
func NewBothLogger(filePath string) (*ZapCtxLogger, error) {
	config := DefaultBothLoggerConfig(filePath)
	return NewZapCtxLoggerWithConfig(config)
}

// NewLoggerWithLevel 创建指定级别的 logger
func NewLoggerWithLevel(mode OutputMode, filePath, level string) (*ZapCtxLogger, error) {
	config := DefaultLoggerConfig()
	config.Output = &OutputConfig{
		Mode:          mode,
		FilePath:      filePath,
		AutoCreateDir: true,
		Level:         level,
	}
	return NewZapCtxLoggerWithConfig(config)
}

// WithContext 设置 context 并返回新的 logger 实例
func (l *ZapCtxLogger) WithContext(ctx context.Context) Logger {
	newLogger := &ZapCtxLogger{
		logger:     l.logger,
		withFields: make([]logger.Field, len(l.withFields)),
		config:     l.config,
		ctx:        ctx,
	}
	copy(newLogger.withFields, l.withFields)
	return newLogger
}

// WithFields 添加字段并返回新的 logger 实例
func (l *ZapCtxLogger) WithFields(args ...logger.Field) Logger {
	newLogger := &ZapCtxLogger{
		logger:     l.logger,
		withFields: make([]logger.Field, len(l.withFields)+len(args)),
		config:     l.config,
		ctx:        l.ctx,
	}
	copy(newLogger.withFields, l.withFields)
	copy(newLogger.withFields[len(l.withFields):], args)
	return newLogger
}

// extractContextFields 从 context 中提取字段
func (l *ZapCtxLogger) extractContextFields(ctx context.Context) []logger.Field {
	if !l.config.EnableContextExtraction || ctx == nil {
		return nil
	}

	if l.config.ContextExtractor != nil {
		return l.config.ContextExtractor.ExtractFields(ctx)
	}

	return nil
}

// prepareFields 准备所有字段（包括 context 字段）
func (l *ZapCtxLogger) prepareFields(ctx context.Context, args []logger.Field) []logger.Field {
	allFields := make([]logger.Field, 0, len(l.withFields)+len(args)+10) // 预留 contex t字段空间

	// 添加预设字段
	allFields = append(allFields, l.withFields...)

	// 添加传入的字段
	allFields = append(allFields, args...)

	// 添加 context 字段
	if ctx != nil {
		contextFields := l.extractContextFields(ctx)
		allFields = append(allFields, contextFields...)
	} else if l.ctx != nil {
		contextFields := l.extractContextFields(l.ctx)
		allFields = append(allFields, contextFields...)
	}

	return allFields
}

// toArgs 将 Field 切片转换为 zap.Field 切片
func (l *ZapCtxLogger) toArgs(args []logger.Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Val))
	}
	return res
}

// 原有的日志方法（向后兼容）
func (l *ZapCtxLogger) Debug(msg string, args ...logger.Field) {
	allFields := l.prepareFields(l.ctx, args)
	l.logger.Debug(msg, l.toArgs(allFields)...)
}

func (l *ZapCtxLogger) Info(msg string, args ...logger.Field) {
	allFields := l.prepareFields(l.ctx, args)
	l.logger.Info(msg, l.toArgs(allFields)...)
}

func (l *ZapCtxLogger) Warn(msg string, args ...logger.Field) {
	allFields := l.prepareFields(l.ctx, args)
	l.logger.Warn(msg, l.toArgs(allFields)...)
}

func (l *ZapCtxLogger) Error(msg string, args ...logger.Field) {
	allFields := l.prepareFields(l.ctx, args)
	l.logger.Error(msg, l.toArgs(allFields)...)
}

// 新增的context支持方法
func (l *ZapCtxLogger) DebugContext(ctx context.Context, msg string, args ...logger.Field) {
	allFields := l.prepareFields(ctx, args)
	l.logger.Debug(msg, l.toArgs(allFields)...)
}

func (l *ZapCtxLogger) InfoContext(ctx context.Context, msg string, args ...logger.Field) {
	allFields := l.prepareFields(ctx, args)
	l.logger.Info(msg, l.toArgs(allFields)...)
}

func (l *ZapCtxLogger) WarnContext(ctx context.Context, msg string, args ...logger.Field) {
	allFields := l.prepareFields(ctx, args)
	l.logger.Warn(msg, l.toArgs(allFields)...)
}

func (l *ZapCtxLogger) ErrorContext(ctx context.Context, msg string, args ...logger.Field) {
	allFields := l.prepareFields(ctx, args)
	l.logger.Error(msg, l.toArgs(allFields)...)
}

// With 方法保持向后兼容（返回 ZapCtxLogger 而不是 Logger 接口）
func (l *ZapCtxLogger) With(args ...logger.Field) *ZapCtxLogger {
	return l.WithFields(args...).(*ZapCtxLogger)
}
