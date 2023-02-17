package main

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	ShellSuffix    = ".cmd"
	IconSuffix     = ".ico"
	TemplateSuffix = ".tmp"
	TemplateDir    = "template"
	AddTemplate    = "toolboxAdd.reg.tmp"
	RemoveTemplate = "toolboxRemove.reg.tmp"
	OutPutDir      = "AutoToolBox"
)

var (
	ErrIdeNotFound = errors.New("no ide script found in the target directory")
)

//go:embed template
var TemplateFs embed.FS

// 输出目录
var targetDir string

// Icon目录
var iconDir string

type JetBrainItemGroup map[string]JetBrainItem

// JetBrainItem
// @Date 2023-02-17 14:57:57
// @Description: JetBrain结构体，代表着一个IDE
type JetBrainItem struct {
	Display   string
	Name      string
	ShellPath string
	IconPath  string
	HKey      string
}

// buildHKey
// @Date 2023-02-17 16:43:38
// @Param flag string
// @Param name string
// @Return string
// @Description: 构建Hkey
func buildHKey(flag string, name string) string {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(flag)
	buffer.WriteByte('.')
	buffer.WriteString(name)
	return buffer.String()
}

// JetBrainToolBox
// @Date 2023-02-17 15:23:01
// @Description: 代表着ToolBox
type JetBrainToolBox struct {
	FLag        string
	Icon        string
	SubCommands string
	IdeGroup    JetBrainItemGroup
}

func NewToolBox() JetBrainToolBox {
	return JetBrainToolBox{
		FLag: "JetBrain",
	}
}

// SetOutPut
// @Date 2023-02-17 14:34:39
// @Param path string
// @Description: 设置输出目录
func SetOutPut(path string) {
	targetDir = path
	iconDir = filepath.Join(targetDir, "ico")
}

// Generate
// @Date 2023-02-17 14:37:11
// @Return error
// @Description: 生成reg注册表脚本
func Generate() error {
	toolBox := NewToolBox()
	// 解析脚本
	group, err := parseShell(toolBox.FLag)
	if err != nil {
		return err
	} else if len(group) == 0 {
		return ErrIdeNotFound
	}
	toolBox.IdeGroup = group
	// 解析SubCommands
	toolBox.SubCommands = parseSubCommands(toolBox.IdeGroup)
	err = parseIdeIcon(&toolBox)
	// 如果icon不存在也不影响使用，不需要返回错误
	if os.IsNotExist(err) {
		Error("the icon directory does not exist")
	}
	// 解析模板文件
	err = parseTemplate(toolBox)
	if err != nil {
		Error("template parse failed", err)
		return err
	}
	return nil
}

// parseShell
// @Date 2023-02-17 14:37:48
// @Return error
// @Method
// @Description: 扫描输出目录下的脚本
func parseShell(flag string) (JetBrainItemGroup, error) {
	dir, err := os.ReadDir(targetDir)
	// 打开目录
	if err != nil {
		return nil, err
	}
	// 初始化map
	jetBrainMap := make(map[string]JetBrainItem, 20)
	// 扫描脚本
	for _, entry := range dir {
		// 如果是cmd脚本的话
		if !entry.IsDir() && strings.Contains(entry.Name(), ShellSuffix) {
			name := entry.Name()[:len(entry.Name())-len(ShellSuffix)]
			display := ToFirstLetterUpper(name)
			// 这里的分隔符必须要进行转义
			shellPath := EscapeRegxPath(filepath.Join(targetDir, entry.Name()))
			hKey := buildHKey(flag, name)
			// 构建结构体放入map
			jetBrainMap[name] = JetBrainItem{
				Display:   display,
				Name:      name,
				ShellPath: shellPath,
				HKey:      hKey,
			}
			Info("IDE ", name, " script detected ", filepath.Join(targetDir, entry.Name()))
		}
	}
	return jetBrainMap, nil
}

// parseSubCommands
// @Date 2023-02-17 15:27:26
// @Return string
// @Return error
// @Description: 解析subcommand
func parseSubCommands(group JetBrainItemGroup) string {
	buffer := bytes.NewBuffer(nil)
	for _, item := range group {
		buffer.WriteString(item.HKey)
		buffer.WriteByte(';')
	}
	return buffer.String()
}

// parseIdeIcon
// @Date 2023-02-17 15:37:03
// @Method
// @Description: 解析Ide图标
func parseIdeIcon(toolbox *JetBrainToolBox) error {
	// 读取Icon目录
	dir, err := os.ReadDir(iconDir)
	// 如果不存在
	if os.IsNotExist(err) {
		return err
	}
	// 遍历目录读取icon文件
	for _, entry := range dir {
		// 不为目录，且是icon文件
		if !entry.IsDir() && strings.Contains(entry.Name(), IconSuffix) {
			name := entry.Name()[:len(entry.Name())-len(IconSuffix)]
			iconPath := filepath.Join(iconDir, entry.Name())
			Info("ico file ", name, " detected ", iconPath)
			// 如果存在对应的ide
			if item, e := toolbox.IdeGroup[name]; e {
				item.IconPath = EscapeRegxPath(iconPath)
				toolbox.IdeGroup[name] = item
			} else if name == "toolbox" {
				toolbox.Icon = EscapeRegxPath(filepath.Join(iconDir, entry.Name()))
			}
		}
	}
	return nil
}

// parseTemplate
// @Date 2023-02-17 16:43:15
// @Param toolbox JetBrainToolBox
// @Description: 解析模板文件
func parseTemplate(toolbox JetBrainToolBox) error {
	// 解析模板文件
	addTem, err := template.ParseFS(TemplateFs, path.Join(TemplateDir, AddTemplate))
	if err != nil {
		return err
	}
	removeTem, err := template.ParseFS(TemplateFs, path.Join(TemplateDir, RemoveTemplate))
	if err != nil {
		return err
	}

	// 创建目录
	dir := filepath.Join(targetDir, OutPutDir)
	if err := Mkdir(dir); err != nil {
		return err
	}

	// 拼接目标输出路径
	addPath := filepath.Join(dir, AddTemplate[:len(AddTemplate)-len(TemplateSuffix)])
	removePath := filepath.Join(dir, RemoveTemplate[:len(RemoveTemplate)-len(TemplateSuffix)])

	// 打开文件
	addRegFile, err := openFile(addPath)
	if err != nil {
		return err
	}
	removeRegFile, err := openFile(removePath)
	if err != nil {
		return err
	}

	// 执行模板解析
	if err := addTem.Execute(addRegFile, toolbox); err != nil {
		return err
	}
	if err := removeTem.Execute(removeRegFile, toolbox); err != nil {
		return err
	}
	Success("reg files has been successfully generated in the directory ", filepath.Join(targetDir, OutPutDir))
	return nil
}

// Exit
// @Date 2023-02-17 16:40:55
// @Method
// @Description: 阻塞或退出
func Exit() {
	Info("press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
