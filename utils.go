package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

const (
	ColorRed   = "\x1b[31m"
	ColorGreen = "\x1b[32m"
	ColorClean = "\x1b[0m"
)

// ToFirstLetterUpper
// @Date 2023-02-17 17:54:43
// @Param str string
// @Return string
// @Description: 第一个字母大写
func ToFirstLetterUpper(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// Mkdir
// @Date 2023-02-17 17:55:04
// @Param dir string
// @Return error
// @Description: 创建一个目录
func Mkdir(dir string) error {
	// 先尝试打开目录
	if _, err := os.Open(dir); os.IsNotExist(err) {
		// 如果文件不存在则创建目录
		if err := os.Mkdir(dir, os.ModeDir); err != nil {
			return err
		} else {
			return nil
		}
	} else if err != nil {
		return err
	} else {
		return nil
	}
}

// openFile
// @Date 2023-02-17 17:55:15
// @Param path string
// @Return *os.File
// @Return error
// @Description: 打开一个文件
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

// EscapeRegxPath
// @Date 2023-02-17 19:07:53
// @Param path string
// @Return string
// @Description: regx脚本转义
func EscapeRegxPath(path string) string {
	return strings.ReplaceAll(path, `\`, `\\`)
}

func Info(args ...any) {
	color.Green("[INFO]\t%-100s", fmt.Sprint(args...))
}

func Error(args ...any) {
	color.Red("[ERROR]\t%-100s", fmt.Sprint(args...))
}

func Success(args ...any) {
	color.Cyan("[TIP]\t%-100s", fmt.Sprint(args...))
}
