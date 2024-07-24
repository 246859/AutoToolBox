package toolbox

import (
	"cmp"
	"encoding/json"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	_StateFile   = "state.json"
	_ToolBoxPath = "/AppData/Local/JetBrains/Toolbox"
)

var _ToolBoxDir string

func NewToolBoxMenuCommand(version string) (*cobra.Command, error) {
	tbmCmd := &cobra.Command{
		Use:          "tbm",
		Version:      version,
		SilenceUsage: true,
		Short:        `ToolBox Menu helper`,
		Long: `Toolbox is located at $HOME/AppData/Local/JetBrains/ by default, unless you 
have modified this location, in most cases you do not need to specify the path.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	// get default installation dir
	defaultTbDIr, err := getDefaultTbDIr()
	if err != nil {
		return nil, err
	}
	tbmCmd.SetVersionTemplate("{{ .Version }}")
	tbmCmd.PersistentFlags().StringVarP(&_ToolBoxDir, "dir", "d", defaultTbDIr, "specify the directory where ToolBox installed")
	tbmCmd.AddCommand(versionCmd)
	tbmCmd.AddCommand(listCmd)
	tbmCmd.AddCommand(addMenuCmd)
	tbmCmd.AddCommand(removeCmd)

	return tbmCmd, nil
}

// ToolBox is a struct to hold the toolbox state.
type ToolBox struct {
	Version string `json:"AppVersion"`
	Tools   []Tool `json:"tools"`
}

// Tool represents an IDE in ToolBox.
type Tool struct {
	Id          string `json:"toolId"`
	Tag         string `json:"tag"`
	Name        string `json:"displayName"`
	Version     string `json:"displayVersion"`
	BuildNumber string `json:"buildNumber"`
	Location    string `json:"installLocation"`
	Command     string `json:"launchCommand"`
}

// GetToolBoxState returns content of ToolBox/state.json
func GetToolBoxState(dir string) (*ToolBox, error) {
	stateFilepath := filepath.Join(dir, _StateFile)
	stateFile, err := os.Open(stateFilepath)
	if err != nil {
		return nil, err
	}
	var toolbox ToolBox
	if err := json.NewDecoder(stateFile).Decode(&toolbox); err != nil {
		return nil, err
	}
	slices.SortFunc(toolbox.Tools, func(a, b Tool) int {
		// safe check
		if a.Name == "" || b.Name == "" {
			return strings.Compare(a.Tag, b.Tag)
		}

		nameA, nameB := strings.ToLower(a.Name), strings.ToLower(b.Name)

		if nameA[0] != nameB[0] {
			return cmp.Compare(nameA[0], nameB[0])
		} else if nameA[0] == nameB[0] {
			return strings.Compare(nameA, nameB)
		} else if nameA == nameB {
			return strings.Compare(a.BuildNumber, b.BuildNumber)
		}
		return strings.Compare(nameA, nameB)
	})
	return &toolbox, err
}
