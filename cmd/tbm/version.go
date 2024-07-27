package main

import (
	"fmt"
	"github.com/246859/AutoToolBox/v3/toolbox"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print ToolBox version",
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := GetToolBoxVersion(ToolBoxDir)
		if err != nil {
			return err
		}
		fmt.Println("ToolBox", version)
		return nil
	},
}

func GetToolBoxVersion(dir string) (string, error) {
	toolBoxState, err := toolbox.GetToolBoxState(dir)
	if err != nil {
		return "", err
	}
	return toolBoxState.Version, nil
}
