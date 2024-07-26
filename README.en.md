# AutoToolBox

**English**|[简体中文](README.md)

> If you are using the old version, please go to [v2.2.0 · 246859/AutoToolBox (github.com)](https://github.com/246859/AutoToolBox/tree/v2.2.0) to view information, or [update to the new version](#Upgrade)

Update: After six years, JetBrains finally started to try to solve the problem of context menu, but the menu item is hidden in the open method, which only works for files, not for directories and directory backgrounds. This is obviously just a very simple function, but it has not been supported for a long time, so this project still needs to exist.

Original issue link: [TBX-2540 (jetbrains.com)](https://youtrack.jetbrains.com/issue/TBX-2540/Associate-file-extenstions-with-correct-Toolbox-app-or-with-the-Toolbox-itself-so-that-files-can-be-launched-from-Windows) [ TBX-2478 (jetbrains.com)](https://youtrack.jetbrains.com/issue/TBX-2478/Windows-Open-Directory-With-Editor)

## Introduction

This is a very simple command line tool for adding a windows right-click menu to the Toolbox App. It has the following features:

- Updating or rolling back the version will not cause the menu to become invalid (when there are multiple versions of the IDE at the same time, only the latest version will be directed)

- You can set the IDE to be opened with administrator privileges
- No need to manually maintain the registry,
- The display order of the menu is synchronized with that in the Toolbox

Here is the effect diagram

<img alt="Effect diagram" src="image/preview.png" width="400" height="300">

## Installation

If you have a go environment and the version is greater than go1.16, you can use the `go install` method to install it, as shown below

```bash
$ go install github.com/246859/AutoToolBox/v3/cmd/tbm@latest
```

Or download the latest binary file directly in Release.

## Use

Version 3.0 is much easier to use. Although there are a few more commands, they are not used in most cases. The only required path parameter is the installation path of the Toolbox. Generally, the Toolbox is installed in the following path by default.

```
$HOME/AppData/Local/Jetbrains/Toolbox/
```

The tool uses the above path by default and does not require additional parameters. If the installation path is modified, you need to use `-d` to specify it (it is best not to modify the installation path of the Toolbox).

Please make sure that **Generate Shell Script** in the settings is turned on, otherwise the tool will not work properly.

<img alt="shellpath" src="image/shellpath.png" width=500 height=200/>


### Upgrade
If you are an old tool user and want to upgrade to a new version, you can use the old generated 'toolboxRemove.reg' to remove the old registry, and then use the new version as follows.

### Start

>  **The tool requires administrator privileges to run properly**

After installation, execute the following command

```bash
$ tbm set -a
```

You can add all locally installed IDEs to the right-click menu. This is the simplest way to use it. In most cases, only this command will be used.

### Commands

```
add         Add ToolBox IDE to existing context menu
clear       clear all the context menu of Toolbox
list        List installed ToolBox IDEs
remove      Remove ToolBox IDEs from context menu
set         Register ToolBox IDEs to context menu
version     Print ToolBox version
```

The following is a brief description of the general function of each command

#### list

```bash
$ tbm list -h
Examples:
  tbm list -c
  tbm list --menu
  tbm list -c --menu

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

The `set` command indicates which IDEs are set as menu items. It will directly overwrite the existing menus. The display order of the menus is the same as in the Toolbox interface.

The simplest way to use it is to directly set all IDEs as menu items. If the number of local IDEs exceeds 16, only the first 16 will be added. This is because the maximum limit of Windows menu items is 16.

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

When registering the menu, you can use `--top` to make the Toolbox menu at the top position

```bash
$ tbm set -a --admin --top
```

<br/>

It should be noted that some products do not provide a stable shell script path or the location of the `exe` file. The following are some of them

```
dotMemory Portable
dotPeek Portable
dotTrace Portable
ReSharper Tools
```

Although they can be added to the menu at this stage, their file structure is not as organized as other IDEs. The `list` command will show which tools are not supported yet, as follows

```bash
$ tbm list
Android Studio                  Koala 2024.1.1 Patch 1
Aqua                            2024.1.2
CLion                           2024.1.4
DataGrip                        2024.1.4
DataSpell                       2024.1.3
dotMemory Portable              2024.1.4                unavailable
dotPeek Portable                2024.1.4                unavailable
dotTrace Portable               2024.1.4                unavailable
Fleet                           1.37.84 Public Preview
```

For them, they will not be added to the menu for the time being.

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

The difference between the `add` command is that it will add new menu items to the existing menu instead of directly overwriting like `set`, and the usage is generally the same.

```bash
$ tbm add GoLand WebStorm
```

However, it does not support `-a`, so all IDEs cannot be added at once.

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



#### clear

```bash
$ tbm clear
clear all the context menu of Toolbox

Usage:
tbm clear [flags]

Flags:
-h, --help help for clear
```

The command `clear` will directly clear all menu items related to Toolbox, including the top-level menu, and will not produce any output. If you do not want to use this tool anymore, you can use this command to clear all registry entries.



## Contribution

1. Fork this repository to your account
2. Create a new branch in the forked repository
3. Submit code changes in the new branch
4. Then initiate a Pull Request to this repository
5. Waiting for Pull Request
