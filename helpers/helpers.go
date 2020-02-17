package helpers

import (
	"os"
	"path/filepath"
)

//GetExePath gets the wroking directory of the EXE
func GetExePath() (path string) {
	exe, _ := os.Executable()
	exPath := filepath.Dir(exe)
	return exPath
}

//FileExists Checkes if a file exists based on file name.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
