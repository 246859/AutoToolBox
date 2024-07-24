package toolbox

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listcount bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed ToolBox IDEs",
	RunE: func(cmd *cobra.Command, args []string) error {
		tools, err := ListLocalTools(_ToolBoxDir)
		if err != nil {
			return err
		}
		if listcount {
			fmt.Println(len(tools))
		} else {
			for _, tool := range tools {
				fmt.Println(tool)
			}
		}
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listcount, "count", "c", false, "count the number of installed tools")
}

// ListLocalTools return local tool list description
func ListLocalTools(dir string) ([]string, error) {
	toolbox, err := GetToolBoxState(dir)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, tool := range toolbox.Tools {
		list = append(list, fmt.Sprintf("%-30s\t%-20s", tool.Name, tool.Version))
	}
	return list, nil
}
