package toolbox

import (
	"cmp"
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"slices"
	"strings"
	"unsafe"
)

// convert a byte slice to a string without allocating new memory.
func bytes2string(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

// convert a string to a byte slice without allocating new memory.
func string2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// compare two version strings.
func compareVersion(version1 string, version2 string) int {
	var i, j int
	for i < len(version1) || j < len(version2) {
		var a, b int
		for ; i < len(version1) && version1[i] != '.'; i++ {
			a = a*10 + int(version1[i]-'0')
		}
		for ; j < len(version2) && version2[j] != '.'; j++ {
			b = b*10 + int(version2[j]-'0')
		}
		if a > b {
			return 1
		} else if a < b {
			return -1
		}

		i++
		j++
	}
	return 0
}

// compare tool name.
func compareName(a string, b string) int {
	// safe check
	nameA, nameB := strings.ToLower(a), strings.ToLower(b)
	if nameA[0] != nameB[0] {
		return cmp.Compare(nameA[0], nameB[0])
	} else if nameA[0] == nameB[0] && nameA != nameB {
		return strings.Compare(nameA, nameB)
	}
	return strings.Compare(nameA, nameB)
}

// open a registry key, create it if it doesn't exist.
func createAndOpen(key registry.Key, path string, access uint32) (registry.Key, error) {
	newk, existing, err := registry.CreateKey(key, path, access)
	if err != nil {
		return 0, err
	}
	if existing {
		return registry.OpenKey(key, path, access)
	}
	return newk, nil
}

// delete a registry key and its sub keys.
func deleteKey(key registry.Key, path string) error {
	parentKey, err := registry.OpenKey(key, path, registry.READ)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	defer parentKey.Close()

	// join path
	if len(path) > 0 && path[len(path)-1] != '\\' {
		path = path + `\`
	}

	subKeyNames, err := parentKey.ReadSubKeyNames(-1)
	if err != nil {
		return err
	}

	for _, name := range subKeyNames {
		subKeyPath := path + name
		err := deleteKey(key, subKeyPath)
		if err != nil {
			return err
		}
	}

	return registry.DeleteKey(key, path)
}

// return command script for the given cmd
func commandScript(cmd string, admin bool) string {
	if admin {
		return fmt.Sprintf(_VBScript, cmd)
	} else {
		return fmt.Sprintf(_CmdScript, cmd)
	}
}

const (
	_SortNames = 0
	_SortOrder = 1
)

func sortTools(tools []*Tool, sortType int) {
	switch sortType {
	case _SortNames:
		slices.SortFunc(tools, func(a, b *Tool) int {
			return compareName(a.Name, b.Name)
		})
	case _SortOrder:
		slices.SortFunc(tools, func(a, b *Tool) int {
			return cmp.Compare(a.order, b.order)
		})
	}
}
