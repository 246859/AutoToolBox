package main

import (
	"fmt"
	"github.com/246859/AutoToolBox/v3/toolbox"
	"github.com/spf13/cobra"
	"os"
)

var Version string

var ToolBoxDir string

var rootCmd = &cobra.Command{
	Use:          "tbm",
	Version:      Version,
	SilenceUsage: true,
	Short:        `ToolBox Menu helper`,
	Long: `tbm is a helper tool to manage ToolBox IDEs context menu on Windows.

Toolbox App is located at $HOME/AppData/Local/JetBrains/ by default, in most cases 
you do not need to specify the path unless you have moved this location. If you did 
do that, use -d to specify the directory.

see more information at https://github.com/246859/AutoToolBox 
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	defaultTbDIr, err := toolbox.DefaultToolboxDir()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	rootCmd.SetVersionTemplate("{{ .Version }}")
	rootCmd.PersistentFlags().StringVar(&ToolBoxDir, "dir", defaultTbDIr, "specify the directory where ToolBox installed")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(clearCmd)
}

func main() {
	rootCmd.Execute()
}
