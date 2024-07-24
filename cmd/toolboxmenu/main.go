package main

import (
	"github.com/246859/AutoToolBox/v3/toolbox"
	"os"
)

var Version string

func main() {
	command, err := toolbox.NewToolBoxMenuCommand(Version)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	} else {
		command.Execute()
	}
}
