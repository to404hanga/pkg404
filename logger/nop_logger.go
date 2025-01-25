package logger

type NopLogger struct {
}

func NewNopLogger() Logger {
	return &NopLogger{}
}

var _ Logger = (*NopLogger)(nil)

func (l *NopLogger) Debug(msg string, args ...Field) {
}

func (l *NopLogger) Info(msg string, args ...Field) {
}

func (l *NopLogger) Warn(msg string, args ...Field) {
}

func (l *NopLogger) Error(msg string, args ...Field) {
}
