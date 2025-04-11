package utils

import (
	"os"
	"path/filepath"
)

type PathUtil struct {
}

func (p *PathUtil) Join(path ...string) string {
	return filepath.Join(path...)
}

func FileExists(path ...string) bool {
	filePath := filepath.Join(path...)
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}
