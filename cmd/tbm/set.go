package main

import (
	"fmt"
	"github.com/246859/AutoToolBox/v3/toolbox"
	"github.com/spf13/cobra"
	"slices"
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
		tools, err := RunSet(ToolBoxDir, args, top, admin, all, update)
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

// RunSet register the specified IDEs to the context menu and return which tools are added successfully.

func RunSet(dir string, targets []string, top, admin, all, update bool) ([]*toolbox.Tool, error) {

	// get all latest tools
	toolboxState, err := toolbox.GetLatestTools(dir, toolbox.SortOrder)
	if err != nil {
		return nil, err
	}
	// collect tools
	tools := toolbox.FindTargetTools(toolboxState.Tools, targets, all || update)

	// select items from menu
	itemsInMenu, _, err := toolbox.ReadSubCommands()
	if err != nil {
		return nil, err
	}
	var temp []*toolbox.Tool
	if update {
		for _, tool := range tools {
			if slices.ContainsFunc(itemsInMenu, func(id string) bool { return tool.Id == id }) {
				temp = append(temp, tool)
			}
		}
		tools = temp
	}

	if len(tools) > toolbox.EntryLimit {
		fmt.Println("Warning: too many tools, only first 16 will be added to the context menu")
		tools = tools[:toolbox.EntryLimit]
	}

	var items []string
	// add menu item
	for _, tool := range tools {
		err := toolbox.SetItem(tool, admin)
		if err != nil {
			return nil, err
		}
		items = append(items, tool.Id)
	}

	// set menu
	if err := toolbox.SetMenu(dir, items, top); err != nil {
		return nil, err
	}

	return tools, nil
}
