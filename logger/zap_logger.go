package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger     *zap.Logger
	withFields []Field
}

var _ Logger = (*ZapLogger)(nil)

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: logger}
}

func (l *ZapLogger) With(args ...Field) *ZapLogger {
	nl := NewZapLogger(l.logger)
	nl.withFields = args
	return nl
}

func (l *ZapLogger) Info(msg string, args ...Field) {
	l.withFields = append(l.withFields, args...)
	l.logger.Info(msg, l.toArgs(l.withFields)...)
}

func (l *ZapLogger) Debug(msg string, args ...Field) {
	l.withFields = append(l.withFields, args...)
	l.logger.Debug(msg, l.toArgs(l.withFields)...)
}

func (l *ZapLogger) Warn(msg string, args ...Field) {
	l.withFields = append(l.withFields, args...)
	l.logger.Warn(msg, l.toArgs(l.withFields)...)
}

func (l *ZapLogger) Error(msg string, args ...Field) {
	l.withFields = append(l.withFields, args...)
	l.logger.Error(msg, l.toArgs(l.withFields)...)
}

func (l *ZapLogger) toArgs(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Val))
	}
	return res
}
