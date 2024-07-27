package toolbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	AppName         = "Toolbox"
	_StateFile      = "state.json"
	_SettingFile    = ".settings.json"
	_ToolBoxPath    = "/AppData/Local/JetBrains/Toolbox"
	_ToolBoxCommand = "/bin/jetbrains-toolbox.exe"
	_ScriptExt      = "cmd"

	// registry keys
	DirectoryBackgroundShell = `Directory\Background\shell\`
	DirectoryShell           = `Directory\shell\`
	CommandStoreShell        = `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\CommandStore\shell\`

	// EntryLimit max limit for cascade menu
	EntryLimit = 16
)

const (
	// using the VBScript to open IDE as admin
	_VBScript  = `mshta vbscript:createobject("shell.application").shellexecute("%s","%%V","","runas",1)(close)`
	_CmdScript = `"%s" "%%V"`
)

// DefaultToolboxDir returns default toolbox installation directory.
func DefaultToolboxDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, _ToolBoxPath), nil
}

// ToolBox is a struct to hold the toolbox state.
type ToolBox struct {
	Version string  `json:"AppVersion"`
	Tools   []*Tool `json:"tools"`

	ShellPath string
}

type Availability int

const (
	available Availability = iota
	legacy
	unavailable
)

func (a Availability) String() string {
	switch a {
	case available:
		return "available"
	case legacy:
		return "legacy"
	case unavailable:
		return "unavailable"
	default:
		return "available"
	}
}

func toolFilter(tools []*Tool, maxAvl Availability) []*Tool {
	var filtered []*Tool
	for _, tool := range tools {
		if tool.Availability <= maxAvl {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// Tool represents an IDE in ToolBox.
type Tool struct {
	Id          string `json:"toolId"`
	Tag         string `json:"tag"`
	Name        string `json:"displayName"`
	Version     string `json:"displayVersion"`
	BuildNumber string `json:"buildNumber"`
	Channel     string `json:"channelId"`
	Location    string `json:"installLocation"`
	// exe
	Command string `json:"launchCommand"`
	// script file
	Script       string
	Availability Availability

	order int
}

// GetToolBoxState returns content of ToolBox/state.json
func GetToolBoxState(dir string) (*ToolBox, error) {

	// get toolbox tools information from state.json
	stateFilepath := filepath.Join(dir, _StateFile)
	stateFile, err := os.Open(stateFilepath)
	if err != nil {
		return nil, err
	}
	var toolbox ToolBox
	if err := json.NewDecoder(stateFile).Decode(&toolbox); err != nil {
		return nil, err
	}

	// read shell path
	settingBytes, err := os.ReadFile(filepath.Join(dir, _SettingFile))
	if err != nil {
		return nil, err
	}
	toolbox.ShellPath = gjson.GetBytes(settingBytes, "shell_scripts.location").String()
	if toolbox.ShellPath == "" {
		defaultDir, err := DefaultToolboxDir()
		if err != nil {
			return nil, err
		}
		toolbox.ShellPath = filepath.Join(defaultDir, "script")
	}

	// get shell script name for per tool from channel file
	for i, tool := range toolbox.Tools {
		channelBytes, err := os.ReadFile(filepath.Join(dir, "channels", fmt.Sprintf("%s.json", tool.Channel)))
		if err != nil {
			return nil, err
		}

		// find shel script from extension
		extensions := gjson.GetBytes(channelBytes, "tool.extensions")
		if extensions.Exists() {
			for _, ext := range extensions.Array() {
				if ext.Get("type").String() == "shell" {
					scriptName := ext.Get("name").String()
					tool.Script = filepath.Join(toolbox.ShellPath, fmt.Sprintf("%s.%s", scriptName, _ScriptExt))
				}
			}
		}

		// judge availability
		if tool.Script == "" && tool.Command == "" {
			tool.Availability = unavailable
		} else if tool.Script == "" {
			tool.Availability = legacy
		}

		// by default, the ordering of tools is depending on state.json that is maintained by Toolbox.
		// Toolbox app will update state.json after you change the order of tools.
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
	SortTools(latestTools, sortType)
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

// FindTargetTools returns tools with specific names
func FindTargetTools(tools []*Tool, targets []string, all bool) []*Tool {
	var targetTools []*Tool
	if all {
		targetTools = tools
	} else {
		for _, target := range targets {
			for _, tool := range tools {
				if tool.Name == target {
					targetTools = append(targetTools, tool)
					break
				}
			}
		}
	}
	return toolFilter(targetTools, legacy)
}

// ReadSubCommands returns current menu items
func ReadSubCommands() ([]string, bool, error) {
	bgShellKey, err := registry.OpenKey(registry.CLASSES_ROOT, DirectoryBackgroundShell+AppName, registry.READ)
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

func SetMenu(dir string, items []string, top bool) error {
	toolboxDisplay := fmt.Sprintf("Open %s Here", AppName)
	toolboxCmd := filepath.Join(dir, _ToolBoxCommand)
	subCommands := strings.Join(items, ";")

	// add directory background shell
	err := SetMenuItem(DirectoryBackgroundShell+AppName, toolboxDisplay, toolboxCmd, subCommands, top)
	if err != nil {
		return err
	}

	// add directory shell
	err = SetMenuItem(DirectoryShell+AppName, toolboxDisplay, toolboxCmd, subCommands, top)
	if err != nil {
		return err
	}
	return nil
}

// SetMenuItem add menu to registry
func SetMenuItem(path, display, command, subCommands string, top bool) error {
	key, err := OpenOrCreateKey(registry.CLASSES_ROOT, path, registry.WRITE)
	if err != nil {
		return err
	}
	defer key.Close()

	// set display content
	if err := key.SetStringValue("MUIVerb", display); err != nil {
		return err
	}

	// set icon
	if err := key.SetStringValue("Icon", command); err != nil {
		return err
	}

	// set sub commands
	if err := key.SetStringValue("SubCommands", subCommands); err != nil {
		return err
	}

	// position
	if top {
		err := key.SetStringValue("Position", "Top")
		if err != nil {
			return err
		}
	}

	return nil
}

// SetItem setItem add items to commandStore shell
func SetItem(tool *Tool, admin bool) error {

	regPath := CommandStoreShell + tool.Id
	ico := filepath.Join(tool.Location, tool.Command)
	script := tool.Script
	// special case for MPS
	if tool.Id == "MPS" {
		ico = filepath.Join(tool.Location, "bin/mps.ico")
	}
	if tool.Availability == legacy {
		script = ico
	} else if tool.Availability == unavailable {
		return nil
	}
	// create or open registry key
	key, err := OpenOrCreateKey(registry.LOCAL_MACHINE, regPath, registry.WRITE)
	if err != nil {
		return err
	}

	// default value
	if err := key.SetStringValue("", fmt.Sprintf("Open %s Here", tool.Name)); err != nil {
		return err
	}
	// set icon
	if err := key.SetStringValue("Icon", ico); err != nil {
		return err
	}

	// command sub key
	cmdKey, err := OpenOrCreateKey(registry.LOCAL_MACHINE, regPath+`\command`, registry.WRITE)
	if err != nil {
		return err
	}

	// set command
	if err := cmdKey.SetStringValue("", commandScript(script, admin)); err != nil {
		return err
	}

	return nil
}
