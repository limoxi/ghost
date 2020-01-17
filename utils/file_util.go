package utils

import "os"

// FileExist 检查文件是否存在
func FileExist(fullPath string) bool{
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}