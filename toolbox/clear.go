package toolbox

import (
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/registry"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear all the context menu of Toolbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ClearToolboxMenu(_ToolBoxDir)
	},
}

func ClearToolboxMenu(dir string) error {
	toolbox, err := GetAllTools(dir)
	if err != nil {
		return err
	}

	// delete all sub-keys whatever can be deleted
	for _, tool := range toolbox.Tools {
		err := deleteKey(registry.LOCAL_MACHINE, _CommandStoreShell+tool.Id)
		if err != nil {
			return err
		}
	}

	if err := deleteKey(registry.CLASSES_ROOT, _DirectoryBackgroundShell+_AppName); err != nil {
		return err
	}
	if err := deleteKey(registry.CLASSES_ROOT, _DirectoryShell+_AppName); err != nil {
		return err
	}

	return nil
}
