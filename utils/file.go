package utils

import (
	"io"
	"os"
)

// IsDir 判断给定路径是否为文件夹
func IsDir(filePath string) bool {
	f, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return f.IsDir()
}

// IsFile 判断给定路径是否为文件
func IsFile(filePath string) bool {
	return !IsDir(filePath)
}

// PathExists 判断文件是否存在
func PathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsFolderEmpty 判断文件夹是否为空
func IsFolderEmpty(folderPath string) bool {
	f, err := os.Open(folderPath)
	if err != nil {
		return false
	}

	_, err = f.ReadDir(1)
	if err == io.EOF {
		return true
	}
	return false
}
