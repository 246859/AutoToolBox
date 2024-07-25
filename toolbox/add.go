package toolbox

import (
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
	"slices"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add ToolBox IDE to existing context menu",
	Long: `Command "add" will append items to the existing context menu instead of overwrite them. 

Use "tbm set -h" for more information.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println(`no tools specified, use "tbm list" to show all installed tools, use "tbm add -h" to get help.`)
			return nil
		}
		tools, err := AddToolboxMenu(_ToolBoxDir, args, top, admin)
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

func AddToolboxMenu(dir string, targets []string, admin, top bool) ([]*Tool, error) {
	// get all latest tools
	toolbox, err := GetLatestTools(dir, _SortOrder)
	if err != nil {
		return nil, err
	}
	// collect tools
	appendTools := FindTargetTools(toolbox.Tools, targets, false)

	if len(appendTools) > _EntryLimit {
		fmt.Println("Warning: too many tools, only first 16 will be added to the context menu")
		appendTools = appendTools[:_EntryLimit]
	}

	// read subCommands
	items, _, err := ReadSubCommands()
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
		cmd := filepath.Join(tool.Location, tool.Command)
		ico := cmd
		if tool.Id == "MPS" {
			ico = filepath.Join(tool.Location, "bin/mps.ico")
		}
		err := setItem(tool.Id, tool.Name, ico, cmd, admin)
		if err != nil {
			return nil, err
		}
	}

	// update menu
	if err := setMenu(items, top); err != nil {
		return nil, err
	}
	return appendTools, nil
}
