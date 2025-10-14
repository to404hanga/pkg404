package loggerv2

import (
	"context"

	"github.com/to404hanga/pkg404/logger"
)

type Logger interface {
	// 原有的日志方法
	logger.Logger

	// 新增的context支持方法
	DebugContext(ctx context.Context, msg string, args ...logger.Field)
	InfoContext(ctx context.Context, msg string, args ...logger.Field)
	WarnContext(ctx context.Context, msg string, args ...logger.Field)
	ErrorContext(ctx context.Context, msg string, args ...logger.Field)

	// 链式调用支持
	WithContext(ctx context.Context) Logger
	WithFields(args ...logger.Field) Logger
}
