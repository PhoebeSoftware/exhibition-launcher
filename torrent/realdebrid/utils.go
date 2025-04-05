package realdebrid

import (
	"os"
	"path/filepath"
)

func (client *RealDebridClient) GetDiskSizeOfAllLinks(path string, unrestrictResponses []UnrestrictResponse) (int64, error) {
	var (
		totalSize int64
	)

	for _, unrestrictResponse := range unrestrictResponses {
		size := int64(unrestrictResponse.Filesize)

		// If the file already exists
		// size = size - filesize
		path = filepath.Join(path, unrestrictResponse.Filename)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			return totalSize, err
		}

		defer file.Close()

		fileStat, err := file.Stat()
		if err != nil {
			return totalSize, err
		}
		size -= fileStat.Size()

		totalSize += size
	}

	return totalSize, nil
}
