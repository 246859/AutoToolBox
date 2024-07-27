package main

import (
	"github.com/246859/AutoToolBox/v3/toolbox"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/registry"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear all the context menu of Toolbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClear(ToolBoxDir)
	},
}

func RunClear(dir string) error {
	toolboxState, err := toolbox.GetAllTools(dir)
	if err != nil {
		return err
	}

	// delete all sub-keys whatever can be deleted
	for _, tool := range toolboxState.Tools {
		err := toolbox.DeleteKey(registry.LOCAL_MACHINE, toolbox.CommandStoreShell+tool.Id)
		if err != nil {
			return err
		}
	}

	if err := toolbox.DeleteKey(registry.CLASSES_ROOT, toolbox.DirectoryBackgroundShell+toolbox.AppName); err != nil {
		return err
	}
	if err := toolbox.DeleteKey(registry.CLASSES_ROOT, toolbox.DirectoryShell+toolbox.AppName); err != nil {
		return err
	}

	return nil
}
