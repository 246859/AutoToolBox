package main

import (
	"fmt"
	"github.com/246859/AutoToolBox/v3/toolbox"
	"github.com/spf13/cobra"
	"slices"
)

var (
	showCount  bool
	showInMenu bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed ToolBox IDEs",
	Long: `Command "list" will list all installed ToolBox IDEs.

Examples:
  tbm list -c 
  tbm list --menu
  tbm list -c --menu
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tools, err := ListToolboxTools(ToolBoxDir, showInMenu)
		if err != nil {
			return err
		}
		if showCount {
			fmt.Println(len(tools))
		} else { // show list
			for _, tool := range tools {
				var tips string
				if tool.Availability > 0 {
					tips = tool.Availability.String()
				}
				fmt.Printf("%-30s\t%-20s\t%-10s\n", tool.Name, tool.Version, tips)
			}
		}
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&showCount, "count", "c", false, "count the number of installed tools")
	listCmd.Flags().BoolVar(&showInMenu, "menu", false, "list the tools shown in the context menu")
}

// ListToolboxTools list local tools
func ListToolboxTools(dir string, showInMenu bool) ([]*toolbox.Tool, error) {
	if !showInMenu {
		toolBox, err := toolbox.GetAllTools(dir)
		if err != nil {
			return nil, err
		}
		return toolBox.Tools, err
	}

	toolboxState, err := toolbox.GetLatestTools(dir, toolbox.SortNames)
	if err != nil {
		return nil, err
	}

	items, exist, err := toolbox.ReadSubCommands()
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, nil
	}

	var tools []*toolbox.Tool
	for _, tool := range toolboxState.Tools {
		if slices.ContainsFunc(items, func(id string) bool { return tool.Id == id }) {
			tools = append(tools, tool)
		}
	}
	toolbox.SortTools(tools, toolbox.SortNames)
	return tools, nil
}
