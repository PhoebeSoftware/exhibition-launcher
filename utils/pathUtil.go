package utils

import "path/filepath"

type PathUtil struct {

}

func (p *PathUtil) Join(path ...string) string {
	return filepath.Join(path...)
}



