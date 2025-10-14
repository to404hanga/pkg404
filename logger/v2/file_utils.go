package loggerv2

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// ensureDir 确保目录存在，如果不存在则创建
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "/" {
		return nil
	}
	
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 创建目录，包括所有必要的父目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// buildZapConfig 根据配置构建 zap.Config
func buildZapConfig(output *OutputConfig) (zap.Config, error) {
	var cfg zap.Config
	
	// 根据级别设置配置
	switch output.Level {
	case "prod", "production":
		cfg = zap.NewProductionConfig()
	case "dev", "development":
		cfg = zap.NewDevelopmentConfig()
	case "test":
		cfg = zap.NewDevelopmentConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
	}
	
	// 设置输出路径
	switch output.Mode {
	case OutputModeConsole:
		// 默认输出到控制台，不需要额外设置
		cfg.OutputPaths = []string{"stdout"}
	case OutputModeFile:
		// 仅输出到文件
		if output.FilePath == "" {
			// 如果没有指定文件路径，使用默认路径
			output.FilePath = "./logs/" + time.Now().Format("20060102") + ".log"
		}
		
		// 如果启用自动创建目录，确保目录存在
		if output.AutoCreateDir {
			if err := ensureDir(output.FilePath); err != nil {
				return cfg, err
			}
		}
		
		cfg.OutputPaths = []string{output.FilePath}
	case OutputModeBoth:
		// 同时输出到控制台和文件
		if output.FilePath == "" {
			output.FilePath = "./logs/" + time.Now().Format("20060102") + ".log"
		}
		
		if output.AutoCreateDir {
			if err := ensureDir(output.FilePath); err != nil {
				return cfg, err
			}
		}
		
		cfg.OutputPaths = []string{"stdout", output.FilePath}
	default:
		// 默认输出到控制台
		cfg.OutputPaths = []string{"stdout"}
	}
	
	return cfg, nil
}

// createZapLogger 创建 zap.Logger 实例
func createZapLogger(output *OutputConfig) (*zap.Logger, error) {
	cfg, err := buildZapConfig(output)
	if err != nil {
		return nil, err
	}
	
	// 构建 logger
	logger, err := cfg.Build(
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCallerSkip(2),
	)
	if err != nil {
		return nil, err
	}
	
	return logger, nil
}