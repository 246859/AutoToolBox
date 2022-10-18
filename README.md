# AutoToolBox

JetBrain的ToolBox一直没有自动添加右键菜单的功能，仅根据exe修改注册表十分的繁琐而且版本更新后就不可用了，

此脚本的生成的注册表脚本是基于shell脚本，版本更新后依然可用，并且可以方便快速修改/删除注册表。

此脚本开发于Windows10系统，默认生成的注册表版本为`5.00`，一些细节需要自行修改。



在使用前确保你的输入的目录应当符合下方的结构

```
dir
|
|---ico
|   |
|   |---idea.ico
|   |
|   |---goland.ico
|   |
|   |---toolbox.ico
|
|---idea.cmd
|
|---goland.cmd
```

