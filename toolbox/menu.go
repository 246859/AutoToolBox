package toolbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	_AppName        = "Toolbox"
	_StateFile      = "state.json"
	_ToolBoxPath    = "/AppData/Local/JetBrains/Toolbox"
	_ToolBoxCommand = "/bin/jetbrains-toolbox.exe"
	// registry keys
	_DirectoryBackgroundShell = `Directory\Background\shell\`
	_DirectoryShell           = `Directory\shell\`
	_CommandStoreShell        = `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\CommandStore\shell\`

	// max limit for cascade menu
	_EntryLimit = 16
)

// returns default toolbox installation directory.
func getDefaultTbDIr() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, _ToolBoxPath), nil
}

var _ToolBoxDir string

func NewToolBoxMenuCommand(version string) (*cobra.Command, error) {
	tbmCmd := &cobra.Command{
		Use:          "tbm",
		Version:      version,
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
	// get default installation dir
	defaultTbDIr, err := getDefaultTbDIr()
	if err != nil {
		return nil, err
	}
	tbmCmd.SetVersionTemplate("{{ .Version }}")
	tbmCmd.PersistentFlags().StringVar(&_ToolBoxDir, "dir", defaultTbDIr, "specify the directory where ToolBox installed")
	tbmCmd.AddCommand(versionCmd)
	tbmCmd.AddCommand(listCmd)
	tbmCmd.AddCommand(setCmd)
	tbmCmd.AddCommand(addCmd)
	tbmCmd.AddCommand(removeCmd)

	return tbmCmd, nil
}

// ToolBox is a struct to hold the toolbox state.
type ToolBox struct {
	Version string  `json:"AppVersion"`
	Tools   []*Tool `json:"tools"`
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

	order           int
	multipleVersion bool
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
	for i, tool := range toolbox.Tools {
		tool.order = i
	}
	return &toolbox, err
}

// GetAllTools return local tool list description
func GetAllTools(dir string) (*ToolBox, error) {
	toolbox, err := GetToolBoxState(dir)
	if err != nil {
		return nil, err
	}
	// sort tools
	slices.SortFunc(toolbox.Tools, func(a, b *Tool) int {
		if a.Name == b.Name {
			return -compareVersion(a.BuildNumber, b.BuildNumber)
		}
		return compareName(a.Name, b.Name)
	})

	// find which tools have different versions
	tools := make(map[string][]*Tool)
	for _, tool := range toolbox.Tools {
		tools[tool.Name] = append(tools[tool.Name], tool)
	}

	for _, list := range tools {
		if len(list) > 1 {
			for _, tool := range list {
				tool.multipleVersion = true
			}
		}
	}

	return toolbox, nil
}

// GetLatestTools returns latest tool list
func GetLatestTools(dir string, sortType int) (*ToolBox, error) {

	tools := make(map[string][]*Tool)
	toolbox, err := GetToolBoxState(dir)
	if err != nil {
		return nil, err
	}

	// collect tools
	for _, tool := range toolbox.Tools {
		tools[tool.Name] = append(tools[tool.Name], tool)
	}

	// collect latest tools
	var latestTools []*Tool
	for _, list := range tools {
		switch len(list) {
		case 0:
			continue
		case 1:
			latestTools = append(latestTools, list[0])
		default:
			latestTools = append(latestTools, FindLatestTool(list))
		}
	}
	sortTools(latestTools, sortType)
	toolbox.Tools = latestTools
	return toolbox, nil
}

// FindLatestTool find the latest tool in a list of tools.
func FindLatestTool(tools []*Tool) *Tool {
	maxTool := &Tool{Version: "0"}
	for _, tool := range tools {
		if compareVersion(tool.Version, maxTool.Version) == 1 {
			maxTool = tool
		}
	}
	return maxTool
}

func FindTargetTools(tools []*Tool, targets []string, all bool) []*Tool {
	var targetTools []*Tool
	if all {
		targetTools = slices.DeleteFunc(tools, func(tool *Tool) bool {
			return slices.ContainsFunc(_ExcludeList, func(s string) bool {
				return tool.Name == s
			})
		})
	} else {
		for _, target := range targets {
			var found bool
			for _, tool := range tools {
				if tool.Name == target {
					targetTools = append(targetTools, tool)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Error: tool %s not found\n", target)
			}
		}
	}
	return targetTools
}

// ReadSubCommands returns current menu items
func ReadSubCommands() ([]string, bool, error) {
	bgShellKey, err := registry.OpenKey(registry.CLASSES_ROOT, _DirectoryBackgroundShell+_AppName, registry.READ)
	if errors.Is(err, registry.ErrNotExist) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	defer bgShellKey.Close()

	value, _, err := bgShellKey.GetStringValue("SubCommands")
	if err != nil {
		return nil, true, err
	}
	if value == "" {
		return nil, true, nil
	}
	return strings.Split(value, ";"), true, err
}

var (
	// unstable tools, they will not show in context menu
	_ExcludeList = []string{
		"dotMemory Portable",
		"dotPeek Portable",
		"dotTrace Portable",
		"ReSharper Tools",
	}
)

func checkExclude(tools []*Tool) error {
	for _, tool := range tools {
		if slices.Contains(_ExcludeList, tool.Name) {
			return fmt.Errorf("unsupported tool: %s", tool.Name)
		}
	}
	return nil
}
