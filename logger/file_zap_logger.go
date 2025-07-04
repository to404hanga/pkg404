package logger

import (
	"os"
	"time"

	"github.com/to404hanga/pkg404/gotools/file"
	"go.uber.org/zap"
)

type FileZapLogger struct {
	logger   *ZapLogger
	used     string
	mode     string
	onlyFile bool
}

var _ Logger = (*FileZapLogger)(nil)

func NewFileZapLogger(mode, outputPath string, onlyFile bool) *FileZapLogger {
	var cfg zap.Config
	switch mode {
	case "dev", "test":
		cfg = zap.NewDevelopmentConfig()
	case "prod":
		cfg = zap.NewProductionConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
	}

	if outputPath == "" {
		outputPath = "./logs/" + time.Now().Format("20060102") + ".log"
		exists, err := file.PathExists(outputPath)
		if err != nil {
			panic(err)
		}
		if !exists {
			os.MkdirAll("./logs", 0777)
			_, err = os.Create(outputPath)
			if err != nil {
				panic(err)
			}
		}
	}
	if onlyFile {
		cfg.OutputPaths = []string{outputPath}
	} else {
		cfg.OutputPaths = append(cfg.OutputPaths, outputPath)
	}

	l, err := cfg.Build(
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(2),
	)
	if err != nil {
		panic(err)
	}
	return &FileZapLogger{
		mode:     mode,
		logger:   NewZapLogger(l),
		used:     outputPath,
		onlyFile: onlyFile,
	}
}

type LoggerFunc func(string, ...Field)

func (l *FileZapLogger) hook(fn LoggerFunc, msg string, args ...Field) {
	now := time.Now().Format("20060102")
	newFile := "./logs/" + now + ".log"
	if l.used != newFile {
		exists, err := file.PathExists(newFile)
		if err != nil {
			panic(err)
		}
		if !exists {
			os.MkdirAll("./logs", 0777)
			_, err := os.Create(newFile)
			if err != nil {
				panic(err)
			}
			nl := NewFileZapLogger(l.mode, newFile, l.onlyFile)
			l = nl
		}
	}
	fn(msg, args...)
}

func (l *FileZapLogger) With(args ...Field) *FileZapLogger {
	nl := NewFileZapLogger(l.mode, l.used, l.onlyFile)
	nl.logger.withFields = args
	return nl
}

func (l *FileZapLogger) Info(msg string, args ...Field) {
	l.hook(l.logger.Info, msg, args...)
}

func (l *FileZapLogger) Debug(msg string, args ...Field) {
	l.hook(l.logger.Debug, msg, args...)
}

func (l *FileZapLogger) Warn(msg string, args ...Field) {
	l.hook(l.logger.Warn, msg, args...)
}

func (l *FileZapLogger) Error(msg string, args ...Field) {
	l.hook(l.logger.Error, msg, args...)
}
