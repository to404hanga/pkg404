package zip

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/itnotebooks/zip"
)

// ZipLib 递归压缩文件或目录，支持AES加密
// 加密方式: Standard, AES128, AES192, AES256(默认)
func ZipLib(dst, src string, encrypt bool, password, enc string) error {
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
		if encrypt {
			encryption := getEncryption(enc)
			fh, err = zFileWriter.Encrypt(header, password, encryption)
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

		ret, err := io.Copy(fh, file)
		if err != nil {
			return err
		}

		log.Printf("added: %s, total: %d\n", path, ret)
		return nil
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
