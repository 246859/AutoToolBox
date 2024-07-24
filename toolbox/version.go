package toolbox

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print toolbox version",
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := GetToolBoxVersion(_ToolBoxDir)
		if err != nil {
			return err
		}
		fmt.Println("ToolBox", version)
		return nil
	},
}

func GetToolBoxVersion(dir string) (string, error) {
	toolbox, err := GetToolBoxState(dir)
	if err != nil {
		return "", err
	}
	return toolbox.Version, nil
}
