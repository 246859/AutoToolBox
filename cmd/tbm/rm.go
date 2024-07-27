package main

import (
	"fmt"
	"github.com/246859/AutoToolBox/v3/toolbox"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/registry"
	"slices"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove ToolBox IDEs from context menu",
	Long: `Command "remove" will remove the specified IDEs from the context menu, use "tbm remove -a" to remove all IDEs.

Example:
  tbm rm GoLand
  tbm rm GoLand WebStorm
  tbm rm -a
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !all && len(args) == 0 {
			fmt.Println(`no tools specified, use "tbm list" to show all installed tools, use "tbm rm -h" to get help.`)
			return nil
		}
		tools, err := RemoveTools(ToolBoxDir, args, all)
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
	removeCmd.Flags().BoolVarP(&silence, "silence", "s", false, "silence output")
	removeCmd.Flags().BoolVarP(&all, "all", "a", false, "remove all")
}

func RemoveTools(dir string, targets []string, all bool) ([]*toolbox.Tool, error) {
	toolboxState, err := toolbox.GetLatestTools(dir, toolbox.SortOrder)
	if err != nil {
		return nil, err
	}

	// get local tools
	preparedTools := toolbox.FindTargetTools(toolboxState.Tools, targets, all)

	// read subcommands
	items, exist, err := toolbox.ReadSubCommands()
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, nil
	}

	// find union set between preparedTools and items
	var temp []*toolbox.Tool
	for _, tool := range preparedTools {
		if slices.Contains(items, tool.Id) {
			temp = append(temp, tool)
		}
	}
	preparedTools = temp

	for _, tool := range preparedTools {
		items = slices.DeleteFunc(items, func(s string) bool {
			return tool.Id == s
		})
	}

	// update menu subCommands
	if err := toolbox.SetMenu(dir, items, false); err != nil {
		return nil, err
	}

	// remove menu item
	var removedTools []*toolbox.Tool
	for _, tool := range preparedTools {
		err := toolbox.DeleteKey(registry.LOCAL_MACHINE, toolbox.CommandStoreShell+tool.Id)
		if err != nil {
			return nil, fmt.Errorf("Error deleting registry key %s:  %v\n", toolbox.CommandStoreShell+tool.Id, err)
		}
		removedTools = append(removedTools, tool)
	}

	// remove top level menu if remove all
	if all {
		if err := toolbox.DeleteKey(registry.CLASSES_ROOT, toolbox.DirectoryBackgroundShell+toolbox.AppName); err != nil {
			return nil, err
		}
		if err := toolbox.DeleteKey(registry.CLASSES_ROOT, toolbox.DirectoryShell+toolbox.AppName); err != nil {
			return nil, err
		}
	}
	return removedTools, nil
}
