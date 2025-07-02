package zip

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/itnotebooks/zip"
)

type EncryptConfig struct {
	Password string
	Enc      string
}

// ZipLib 递归压缩文件或目录，支持AES加密
//
// 如果 cfg 为空，则不加密
//
// 加密方式: Standard, AES128, AES192, AES256(默认)
func ZipLib(dst, src string, cfg ...EncryptConfig) error {
	// 创建压缩文件
	zfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zfile.Close()

	// 创建zip写入器
	zFileWriter := zip.NewWriter(zfile)
	defer func() {
		if err := zFileWriter.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 遍历源路径下的所有文件
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		// 跳过zip文件
		if strings.HasSuffix(path, ".zip") {
			return nil
		}

		// 创建文件头
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 设置相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// 目录处理
		if fi.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// 创建文件写入器
		var fh io.Writer
		if len(cfg) != 0 {
			encryption := getEncryption(cfg[0].Enc)
			fh, err = zFileWriter.Encrypt(header, cfg[0].Password, encryption)
		} else {
			fh, err = zFileWriter.CreateHeader(header)
		}
		if err != nil {
			return err
		}

		// 只处理常规文件
		if !header.Mode().IsRegular() {
			return nil
		}

		// 复制文件内容
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(fh, file)
		return err
	})
}

// getEncryption 根据字符串返回对应的加密方式
func getEncryption(enc string) zip.EncryptionMethod {
	switch enc {
	case "Standard":
		return zip.StandardEncryption
	case "AES128":
		return zip.AES128Encryption
	case "AES192":
		return zip.AES192Encryption
	default:
		return zip.AES256Encryption
	}
}

const (
	Standard = "Standard"
	AES128   = "AES128"
	AES192   = "AES192"
	Default  = "AES256"
)

// UnzipLib 解压加密的zip文件
// src: 源zip文件路径
// dst: 目标解压目录路径
// password: 解密密码
func UnzipLib(src, dst, password string) error {
	// 使用支持加密的zip库打开文件
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	os.MkdirAll(dst, 0755)

	for _, file := range zipReader.File {
		targetPath := filepath.Join(dst, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, file.FileInfo().Mode())
			continue
		}

		err = extractEncryptedFile(file, targetPath, password)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractEncryptedFile 提取加密的单个文件
func extractEncryptedFile(file *zip.File, targetPath, password string) error {
	os.MkdirAll(filepath.Dir(targetPath), 0755)

	// 设置密码（如果文件是加密的）
	if file.IsEncrypted() {
		file.SetPassword(password)
	}

	// 打开文件
	zipFile, err := file.Open()
	if err != nil {
		return err
	}
	defer zipFile.Close()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, zipFile)
	return err
}
