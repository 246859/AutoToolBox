# AutoToolBox

**English**|[简体中文](README.md)

> If you are using the old version, please go to [v2.2.0 · 246859/AutoToolBox (github.com)](https://github.com/246859/AutoToolBox/tree/v2.2.0) to view information, or update to the new version

Update: After six years, JetBrains finally started to try to solve the problem of context menu, but the menu item is hidden in the open method, which only works for files, not for directories and directory backgrounds. This is obviously just a very simple function, but it has not been supported for a long time, so this project still needs to exist.

Original issue link: [TBX-2540 (jetbrains.com)](https://youtrack.jetbrains.com/issue/TBX-2540/Associate-file-extenstions-with-correct-Toolbox-app-or-with-the-Toolbox-itself-so-that-files-can-be-launched-from-Windows) [ TBX-2478 (jetbrains.com)](https://youtrack.jetbrains.com/issue/TBX-2478/Windows-Open-Directory-With-Editor)

## Introduction

This is a very simple command line tool used to add a windows right-click menu to the Toolbox App. It has the following features

- Updating or rolling back the version will not be invalid (when there are multiple versions of the IDE at the same time, only the latest version will be directed)

- You can set it to open the IDE with administrator privileges
- No need to manually maintain the registry

The following is the effect picture

<img alt="Effect picture" src="https://public-1308755698.cos.ap-chongqing.myqcloud.com//upload/202407251742834.png" width="400" height="300">

## Installation

If you have a go environment and the version is greater than go1.16, you can use the `go install`  command to install it, as shown below

```bash
$ go install github.com/246859/AutoToolBox/v3/cmd/tbm@latest
```

Or download the latest binary file directly in Release.

## Use

The 3.0 version of the tool is much simpler to use. Although there are a few more commands, they are not used in most cases. The only path parameter required is the installation path of the Toolbox. Generally, the Toolbox is installed by default in the following path.

```
$HOME/AppData/Local/Jetbrains/Toolbox/
```

The tool uses the above path by default, and no additional parameters are required. If the installation path is modified, it needs to be specified with `-d` (it is best not to modify the installation path of the Toolbox).

Please make sure that **Generate Shell Script** in the settings is turned on, otherwise the tool will not work properly.

<img alt="shellpath" src="https://public-1308755698.cos.ap-chongqing.myqcloud.com//upload/202407251742830.png" width=500 height=200/>

### Start

After installation, execute the following command

```bash
$ tbm set -a
```

You can add all locally installed IDEs to the right-click menu. This is the simplest way to use it. In most cases, only this command will be used.

### Commands

```
Available Commands:
  add         Add ToolBox IDE to existing context menu
  list        List installed ToolBox IDEs
  remove      Remove ToolBox IDEs from context menu
  set         Register ToolBox IDEs to context menu
  version     Print ToolBox version
```

The following is a brief description of the general function of each command

#### list

```bash
$ tbm list -h
Usage:
  tbm list [flags]

Flags:
  -c, --count   count the number of installed tools
  -h, --help    help for list
      --menu    list the tools shown in the context menu
```

`list` command is used to view all locally installed IDEs, for example

```bash
$ tbm list
Android Studio                  Koala 2024.1.1 Patch 1
Aqua                            2024.1.2
CLion                           2024.1.4
DataGrip                        2024.1.4
GoLand                          2024.1.4
GoLand                          2023.3.7
IntelliJ IDEA Community Edition 2024.1.4
IntelliJ IDEA Ultimate          2024.1.4
MPS                             2023.3.1
PhpStorm                        2024.1.4
PyCharm Community               2024.1.4
PyCharm Professional            2024.1.4
```

Check the number

```bash
$ tbm list -c
25
```

Check all items added to the menu

```bash
$ tbm list --menu
Aqua                            2024.1.2
CLion                           2024.1.4
DataGrip                        2024.1.4
DataSpell                       2024.1.3
Fleet                           1.37.84 Public Preview
GoLand                          2024.1.4
IntelliJ IDEA Ultimate          2024.1.4
MPS                             2023.3.1
PhpStorm                        2024.1.4
PyCharm Professional            2024.1.4
Rider                           2024.1.4
RubyMine                        2024.1.4
```

View the number of items added to the menu

```bash
$ tbm list --menu -c
16
```

#### set

```bash
$ tbm set -h
Usage:
  tbm set [flags]

Flags:
      --admin     run as admin
  -a, --all       select all
  -h, --help      help for set
  -s, --silence   silence output
      --top       place toolbox menu at top of context menu
  -u, --update    only select current menu items
```

The `set` command is used to register the locally installed IDE to the right-click menu. It works by overwriting. Each execution will overwrite the previous menu. If you want to add them one by one, you can consider using the `add` command.

The simplest use is to add all directly. If the number of local IDEs exceeds 16, only the first 16 will be added. This is because the maximum number of Windows menu items is 16.

```bash
$ tbm set -a
Warning: too many tools, only first 16 will be added to the context menu
GoLand
IntelliJ IDEA Ultimate
PyCharm Professional
WebStorm
RustRover
Aqua
Writerside
Fleet
DataSpell
CLion
PhpStorm
DataGrip
Rider
RubyMine
Space Desktop
MPS
```

Or specify separately

```bash
$ tbm set GoLand WebStorm
```

If you need to run the IDE with administrator privileges, you can add `--admin`, as shown below,

```bash
$ tbm set -a --admin
```

Using `--update` will only update existing menu items, not add new ones. If there are multiple versions of the same IDE, this command can guide it to the latest version.

```bash
$ tbm set --update
```

When registering a menu, you can use `--top` to make the Toolbox menu at the top position

```bash
$ tbm set -a --admin --top
```

#### add

```bash
$ tbm add -h
Usage:
  tbm add [flags]

Flags:
      --admin     run as admin
  -h, --help      help for add
  -s, --silence   silence output
      --top       place toolbox menu at top of context menu
```

The difference between the `add` command is that it will add new menu items to the existing menu, instead of directly overwriting like `set`, and the usage is generally the same.

```bash
$ tbm add GoLand WebStorm
```

However, it does not support `-a`, and cannot add all IDEs at once.

#### rmove

```bash
$ tbm remove -h
Command "remove" will remove the specified IDEs from the context menu, use "tbm remove -a" to remove all IDEs.

Usage:
  tbm remove [flags]

Aliases:
  remove, rm

Flags:
  -a, --all       remove all
  -h, --help      help for remove
  -s, --silence   silence output
```

The `remove` command is used to remove menu items

```bash
$ tbm rm GoLand WebStorm
```

Use `-a` to remove all

```bash
$ tbm rm -a
```

## Contribution

1. Fork this repository to your account
2. Create a new branch in the forked repository
3. Submit code changes in the new branch
4. Then initiate a Pull Request to this repository
5. Waiting for Pull Request
