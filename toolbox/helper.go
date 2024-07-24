package toolbox

import (
	"os"
	"path/filepath"
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

// returns default toolbox installation directory.
func getDefaultTbDIr() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, _ToolBoxPath), nil
}
