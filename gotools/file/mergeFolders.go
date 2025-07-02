package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// CopyMultipleDirsConcurrent 并发拷贝多个目录
func CopyMultipleDirsConcurrent(srcDirs []string, dstDir string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, len(srcDirs))

	// 并发拷贝每个目录
	for _, srcDir := range srcDirs {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()
			if err := copySingleDir(src, dstDir); err != nil {
				errorChan <- fmt.Errorf("拷贝目录 %s 失败: %v", src, err)
			}
		}(srcDir)
	}

	// 等待所有 goroutine 完成
	wg.Wait()
	close(errorChan)

	// 返回第一个错误
	select {
	case err := <-errorChan:
		return err
	default:
		return nil
	}
}

// copySingleDir 拷贝单个目录
func copySingleDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算目标路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		// 处理目录
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 处理文件
		return copyFile(path, dstPath)
	})
}

// copyFile 拷贝单个文件
func copyFile(src, dst string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 拷贝文件内容
	_, err = io.Copy(dstFile, srcFile)
	return err
}
