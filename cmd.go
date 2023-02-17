package main

import (
	"errors"
	"flag"
	"path/filepath"
)

var (
	ErrToolBoxPath = errors.New("invalid file path or not an absolute path")
)

// FlagParse
// @Date 2023-02-17 14:22:42
// @Return string
// @Description: 解析命令行参数
func FlagParse() (string, error) {
	var toolBoxPath string
	flag.StringVar(&toolBoxPath, "path", " default empty value ", "toolbox shell script dir absolute path")
	flag.Parse()
	if len(toolBoxPath) == 0 || !filepath.IsAbs(toolBoxPath) {
		return toolBoxPath, ErrToolBoxPath
	}
	return toolBoxPath, nil
}
