package main

func main() {
	defer Exit()
	// 解析命令参数
	toolboxPath, err := FlagParse()
	if err != nil {
		Error(err, toolboxPath)
		return
	}
	//设置输出目录
	SetOutPut(toolboxPath)
	// 生成reg脚本
	if err := Generate(); err != nil {
		Error(err)
	} else {
		// 没有运行错误才会输出这一句
		Success("If this gadget can help you, please go to Github:https://github.com/246859/AutoToolBox and give the author a star. Thank you very much!")
	}
}
