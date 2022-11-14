# AutoToolBox


JetBrain's ToolBox has not automatically added the function of right-click menu, only according to exe to modify the registry is very cumbersome and not available after the version update.

The resulting registry script from this script is a shell script that will still be available after an updated version and can facilitate quick modification/deletion of the registry.

This script is developed on Windows 10 systems, the default generated registry version is '5.00', some details need to be modified by yourself.

The program does not directly modify the registry, but generates two registry script files in the input directory, through which the registry is modified
'toolboxAdd.reg' is responsible for adding registry information
'toolboxDelete.reg' is responsible for deleting registry information
Note You must enter an absolute path!!!

Make sure your input directory should conform to the structure below before using it

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

As a reminder, the exe copy-paste path is the right mouse button.

JetBrain的ToolBox一直没有自动添加右键菜单的功能，仅根据exe修改注册表十分的繁琐而且版本更新后就不可用了，

此脚本的生成的注册表脚本是基于shell脚本，版本更新后依然可用，并且可以方便快速修改/删除注册表。

此脚本开发于Windows10系统，默认生成的注册表版本为`5.00`，一些细节需要自行修改。

程序并不会直接修改注册表，而是在输入的目录下生成两个注册表脚本文件，通过脚本文件来修改注册表
`toolboxAdd.reg`负责增加注册表信息
`toolboxDelete.reg` 负责删除注册表信息
注意必须输入绝对路径！！！

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

提醒一下exe复制粘贴路径是鼠标右键。


