package toolbox

import (
	"os"
	"strings"
)

// openFile 打开一个文件
func openFile(path string) (*os.File, error) {
	// 先尝试打开一个文件
	if open, err := os.OpenFile(path, os.O_RDWR, os.ModePerm); os.IsNotExist(err) {
		if file, err := os.Create(path); err != nil {
			return nil, err
		} else {
			return file, err
		}
	} else if err != nil {
		return nil, err
	} else {
		return open, nil
	}
}

// EscapeRegxPath regx脚本转义
func EscapeRegxPath(path string) string {
	return strings.ReplaceAll(path, `\`, `\\`)
}
