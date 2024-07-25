package toolbox

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
	"slices"
	"strings"
)

var (
	// whether to set ToolBox at top of context menu
	top bool
	// Run ToolBox IDE as admin mode
	admin bool
	// slice output
	silence bool
	// update subCommands
	update bool
	// select all
	all bool
)

const (
	// using the VBScript to open IDE as admin
	_VBScript  = `mshta vbscript:createobject("shell.application").shellexecute("%s","%%V","","runas",1)(close)`
	_CmdScript = `"%s" "%%V"`
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Register ToolBox IDEs to context menu",
	Long: `Command "set" will create a new menu to overwrite existing menu, it will set all by default.
If you want to append items to the existing menu, use "tbm add".

Microsoft limits the number of items in the context menu to no more than 16, so only the first 16 tools 
will be set to the context menu. If there has different versions of the same IDE, only the latest 
will be set.

The default order of IDEs in the menu is determined by ToolBox/state.json, that is, 
the download time of the installation tool is sorted. In other case, the order of the 
menus depends on the order of args.

Examples: 
  tbm set -a
  tbm set Goland
  tbm set Goland CLion Webstorm
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !all && !update && len(args) == 0 {
			fmt.Println(`no tools specified, use "tbm list" to show all installed tools, use "tbm rm -h" to get help.`)
			return nil
		}
		tools, err := SetToolBoxMenu(_ToolBoxDir, args, top, admin, all, update)
		if err != nil {
			return err
		}
		if !silence {
			for _, tool := range tools {
				fmt.Println(tool.Name)
			}
		}
		return nil
	},
}

func init() {
	setCmd.Flags().BoolVar(&top, "top", false, "place toolbox menu at top of context menu")
	setCmd.Flags().BoolVar(&admin, "admin", false, "run as admin")
	setCmd.Flags().BoolVarP(&silence, "silence", "s", false, "silence output")
	setCmd.Flags().BoolVarP(&update, "update", "u", false, "only select current menu items")
	setCmd.Flags().BoolVarP(&all, "all", "a", false, "select all")
}

// SetToolBoxMenu register the specified IDEs to the context menu and return which tools are added successfully.

func SetToolBoxMenu(dir string, targets []string, top, admin, all, update bool) ([]*Tool, error) {

	// get all latest tools
	toolbox, err := GetLatestTools(dir, _SortOrder)
	if err != nil {
		return nil, err
	}
	// collect tools
	tools := FindTargetTools(toolbox.Tools, targets, all || update)

	// select items from menu
	itemsInMenu, _, err := ReadSubCommands()
	if err != nil {
		return nil, err
	}
	var temp []*Tool
	if update {
		for _, tool := range tools {
			if slices.ContainsFunc(itemsInMenu, func(id string) bool { return tool.Id == id }) {
				temp = append(temp, tool)
			}
		}
		tools = temp
	}

	if len(tools) > _EntryLimit {
		fmt.Println("Warning: too many tools, only first 16 will be added to the context menu")
		tools = tools[:_EntryLimit]
	}

	var items []string
	// add menu item
	for _, tool := range tools {
		err := setItem(tool, admin)
		if err != nil {
			return nil, err
		}
		items = append(items, tool.Id)
	}

	if err := setMenu(items, top); err != nil {
		return nil, err
	}

	return tools, nil
}

func setMenu(items []string, top bool) error {
	toolboxDisplay := fmt.Sprintf("Open %s Here", _AppName)
	toolboxCmd := filepath.Join(_ToolBoxDir, _ToolBoxCommand)
	subCommands := strings.Join(items, ";")

	// add directory background shell
	err := setMenuItem(_DirectoryBackgroundShell+_AppName, toolboxDisplay, toolboxCmd, subCommands, top)
	if err != nil {
		return err
	}

	// add directory shell
	err = setMenuItem(_DirectoryShell+_AppName, toolboxDisplay, toolboxCmd, subCommands, top)
	if err != nil {
		return err
	}
	return nil
}

// setMenuItem add menu to registry
func setMenuItem(path, display, command, subCommands string, top bool) error {
	key, err := createAndOpen(registry.CLASSES_ROOT, path, registry.WRITE)
	if err != nil {
		return err
	}
	defer key.Close()

	// set display content
	if err := key.SetStringValue("MUIVerb", display); err != nil {
		return err
	}

	// set icon
	if err := key.SetStringValue("Icon", command); err != nil {
		return err
	}

	// set sub commands
	if err := key.SetStringValue("SubCommands", subCommands); err != nil {
		return err
	}

	// position
	if top {
		err := key.SetStringValue("Position", "Top")
		if err != nil {
			return err
		}
	}

	return nil
}

// setItem add items to commandStore shell
func setItem(tool *Tool, admin bool) error {

	regPath := _CommandStoreShell + tool.Id
	ico := filepath.Join(tool.Location, tool.Command)
	script := tool.Script
	// special case for MPS
	if tool.Id == "MPS" {
		ico = filepath.Join(tool.Location, "bin/mps.ico")
	}
	if tool.availability == legacy {
		script = ico
	} else if tool.availability == unavailable {
		return nil
	}
	// create or open registry key
	key, err := createAndOpen(registry.LOCAL_MACHINE, regPath, registry.WRITE)
	if err != nil {
		return err
	}

	// default value
	if err := key.SetStringValue("", fmt.Sprintf("Open %s Here", tool.Name)); err != nil {
		return err
	}
	// set icon
	if err := key.SetStringValue("Icon", ico); err != nil {
		return err
	}

	// command sub key
	cmdKey, err := createAndOpen(registry.LOCAL_MACHINE, regPath+`\command`, registry.WRITE)
	if err != nil {
		return err
	}

	// set command
	if err := cmdKey.SetStringValue("", commandScript(script, admin)); err != nil {
		return err
	}

	return nil
}
