package toolbox

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/246859/AutoToolBox/v2/assets"
	"github.com/dstgo/filebox"
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
	OutPutDir      = "regx"
)

var displayNames = map[string]string{
	"goland":     "Goland",
	"idea":       "IntelliJ IDEA",
	"pycharm":    "Pycharm",
	"datagrip":   "DataGrip",
	"clion":      "CLion",
	"webstorm":   "WebStorm",
	"studio":     "Android Studio",
	"fleet":      "Fleet",
	"rustrover":  "RustRover",
	"writerside": "Writerside",
	"dataspell":  "DataSpell",
	"rubymine":   "RubyMine",
	"rider":      "Rider",
	"phpstorm":   "PhpStorm",
}

type JetBrainItemGroup map[string]JetBrainItem

// JetBrainItem JetBrain结构体，代表着一个IDE
type JetBrainItem struct {
	Display   string
	Name      string
	ShellPath string
	IconPath  string
	HKey      string
}

// JetBrainToolBox 代表着ToolBox
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

// Generate 生成reg注册表脚本
func Generate(output string) error {
	toolBox := NewToolBox()
	// 解析脚本
	group, err := parseShell(output, toolBox.FLag)
	if err != nil {
		return err
	} else if len(group) == 0 {
		return errors.New("no ide script found in the target directory")
	}
	toolBox.IdeGroup = group
	// 解析SubCommands
	toolBox.SubCommands = parseSubCommands(toolBox.IdeGroup)
	// 生成图标文件
	iconDir := path.Join(output, "ico")
	if err := filebox.CopyFs(assets.Fs, "ico", output); err != nil {
		return err
	}
	// 如果icon不存在也不影响使用，不需要返回错误
	if err := parseIdeIcon(iconDir, &toolBox); os.IsNotExist(err) {
		fmt.Printf("warning: ico dir not found output %s\n", output)
	} else if err != nil {
		return err
	}
	// 解析模板
	if err := parseTemplate(output, toolBox); err != nil {
		return err
	}
	return nil
}

// parseShell 扫描输出目录下的脚本
func parseShell(input string, flag string) (JetBrainItemGroup, error) {
	dir, err := os.ReadDir(input)
	// 打开目录
	if err != nil {
		return nil, err
	}
	// 初始化map
	jetBrainMap := make(map[string]JetBrainItem, 20)
	var names []string
	// 扫描脚本
	for _, entry := range dir {
		// 如果是cmd脚本的话
		if !entry.IsDir() && strings.Contains(entry.Name(), ShellSuffix) {
			name := strings.ToLower(entry.Name()[:len(entry.Name())-len(ShellSuffix)])
			display := displayNames[name]
			// 这里的分隔符必须要进行转义
			shellPath := EscapeRegxPath(filepath.Join(input, entry.Name()))
			hKey := fmt.Sprintf("%s.%s", flag, name)
			// 构建结构体放入map
			jetBrainMap[name] = JetBrainItem{
				Display:   display,
				Name:      name,
				ShellPath: shellPath,
				HKey:      hKey,
			}
			names = append(names, name)
		}
	}
	fmt.Printf("found %s\n", strings.Join(names, ","))
	return jetBrainMap, nil
}

// parseSubCommands 解析subcommand
func parseSubCommands(group JetBrainItemGroup) string {
	buffer := bytes.NewBuffer(nil)
	for _, item := range group {
		buffer.WriteString(item.HKey)
		buffer.WriteByte(';')
	}
	return buffer.String()
}

// parseIdeIcon 解析Ide图标
func parseIdeIcon(output string, toolbox *JetBrainToolBox) error {
	// 读取Icon目录
	dir, err := os.ReadDir(output)
	// 如果不存在
	if os.IsNotExist(err) {
		return err
	}
	// 遍历目录读取icon文件
	for _, entry := range dir {
		// 不为目录，且是icon文件
		if !entry.IsDir() && strings.Contains(entry.Name(), IconSuffix) {
			name := entry.Name()[:len(entry.Name())-len(IconSuffix)]
			iconPath := filepath.Join(output, entry.Name())
			// 如果存在对应的ide
			if item, e := toolbox.IdeGroup[name]; e {
				item.IconPath = EscapeRegxPath(iconPath)
				toolbox.IdeGroup[name] = item
			} else if name == "toolbox" {
				toolbox.Icon = EscapeRegxPath(filepath.Join(output, entry.Name()))
			}
		}
	}
	return nil
}

// parseTemplate 解析模板文件
func parseTemplate(target string, toolbox JetBrainToolBox) error {
	// 解析模板文件
	addTem, err := template.ParseFS(assets.Fs, path.Join(TemplateDir, AddTemplate))
	if err != nil {
		return err
	}
	removeTem, err := template.ParseFS(assets.Fs, path.Join(TemplateDir, RemoveTemplate))
	if err != nil {
		return err
	}

	// 创建目录
	dir := filepath.Join(target, OutPutDir)
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
	fmt.Println("generated")
	return nil
}
