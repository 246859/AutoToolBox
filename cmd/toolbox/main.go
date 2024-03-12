package main

import (
	"fmt"
	"github.com/246859/AutoToolBox/v2/toolbox"
	"github.com/spf13/cobra"
	"os/user"
)

var in string
var rootCmd = &cobra.Command{
	Use:           "toolbox [command]",
	Short:         "toolbox is a command line tool to generate win regex scripts for jetbrain ide",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 生成reg脚本
		if err := toolbox.Generate(in); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	current, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}
	userHomeDir := fmt.Sprintf(`%s\AppData\Local\JetBrains\Toolbox\scripts`, current.HomeDir)
	rootCmd.Flags().StringVarP(&in, "input", "i", userHomeDir, "input directory")
}

func main() {
	rootCmd.Execute()
}
