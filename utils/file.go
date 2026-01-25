package utils

import (
	"os"
)

// fileExists 判断文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
