package main

import (
	"fmt"
	"github.com/246859/AutoToolBox/v3/toolbox"
	"github.com/spf13/cobra"
	"slices"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add ToolBox IDE to existing context menu",
	Long: `Command "add" will append items to the existing context menu instead of overwrite them. 

Use "tbm set -h" for more information.

Examples:
  tbm add GoLand --admin
  tbm add GoLand WebStorm 
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println(`no tools specified, use "tbm list" to show all installed tools, use "tbm add -h" to get help.`)
			return nil
		}
		tools, err := RunAdd(ToolBoxDir, args, top, admin)
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
	addCmd.Flags().BoolVar(&top, "top", false, "place toolbox menu at top of context menu")
	addCmd.Flags().BoolVar(&admin, "admin", false, "run as admin")
	addCmd.Flags().BoolVarP(&silence, "silence", "s", false, "silence output")
}

func RunAdd(dir string, targets []string, admin, top bool) ([]*toolbox.Tool, error) {
	// get all latest tools
	toolboxState, err := toolbox.GetLatestTools(dir, toolbox.SortOrder)
	if err != nil {
		return nil, err
	}
	// collect tools
	appendTools := toolbox.FindTargetTools(toolboxState.Tools, targets, false)

	if len(appendTools) > toolbox.EntryLimit {
		fmt.Println("Warning: too many tools, only first 16 will be added to the context menu")
		appendTools = appendTools[:toolbox.EntryLimit]
	}

	// read subCommands
	items, _, err := toolbox.ReadSubCommands()
	if err != nil {
		return nil, err
	}

	// collect items to be saved
	for _, tool := range appendTools {
		if !slices.ContainsFunc(items, func(item string) bool { return item == tool.Id }) {
			items = append(items, tool.Id)
		}
	}

	// create new subcommands
	for _, tool := range appendTools {
		err := toolbox.SetItem(tool, all)
		if err != nil {
			return nil, err
		}
	}

	// update menu
	if err := toolbox.SetMenu(dir, items, top); err != nil {
		return nil, err
	}
	return appendTools, nil
}
