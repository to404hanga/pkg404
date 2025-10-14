package loggerv2

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// createFileWriter 创建文件写入器，支持自动创建目录和文件
func createFileWriter(filePath string, autoCreate bool) (io.Writer, error) {
	if !autoCreate {
		// 不自动创建，直接尝试打开文件
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		return file, nil
	}

	// 自动创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// 创建或打开文件
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	return file, nil
}

// createZapCore 根据输出配置创建zap core
func createZapCore(config OutputConfig, development bool) (zapcore.Core, error) {
	// 设置编码器配置
	var encoderConfig zapcore.EncoderConfig
	if development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	// 创建编码器
	var encoder zapcore.Encoder
	if development {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 设置日志级别
	level := zapcore.InfoLevel
	if development {
		level = zapcore.DebugLevel
	}

	// 根据输出类型创建WriteSyncer
	var writeSyncer zapcore.WriteSyncer
	var err error

	switch config.Type {
	case OutputConsole:
		// 输出到控制台
		if config.Writer != nil {
			writeSyncer = zapcore.AddSync(config.Writer)
		} else {
			writeSyncer = zapcore.AddSync(os.Stdout)
		}

	case OutputFile:
		// 输出到文件
		if config.FilePath == "" {
			return nil, fmt.Errorf("file path is required for file output")
		}

		var writer io.Writer
		if config.Writer != nil {
			writer = config.Writer
		} else {
			writer, err = createFileWriter(config.FilePath, config.AutoCreateFile)
			if err != nil {
				return nil, err
			}
		}
		writeSyncer = zapcore.AddSync(writer)

	case OutputBoth:
		// 同时输出到控制台和文件
		if config.FilePath == "" {
			return nil, fmt.Errorf("file path is required for both output")
		}

		var fileWriter io.Writer
		if config.Writer != nil {
			fileWriter = config.Writer
		} else {
			fileWriter, err = createFileWriter(config.FilePath, config.AutoCreateFile)
			if err != nil {
				return nil, err
			}
		}

		// 创建多重写入器
		writeSyncer = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(fileWriter),
		)

	default:
		return nil, fmt.Errorf("unsupported output type: %d", config.Type)
	}

	// 创建core
	core := zapcore.NewCore(encoder, writeSyncer, level)
	return core, nil
}

// createZapLogger 根据配置创建zap logger
func createZapLogger(config LoggerConfig) (*zap.Logger, error) {
	core, err := createZapCore(config.Output, config.Development)
	if err != nil {
		return nil, err
	}

	// 创建logger选项
	var options []zap.Option
	if config.Development {
		options = append(options, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		options = append(options, zap.AddCaller())
	}

	logger := zap.New(core, options...)
	return logger, nil
}
