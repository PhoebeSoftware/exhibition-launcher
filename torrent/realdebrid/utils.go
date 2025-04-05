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
		size := unrestrictResponse.Filesize

		// If the file already exists
		// size = size - filesize
		pathToFile := filepath.Join(path, unrestrictResponse.Filename)
		file, err := os.OpenFile(pathToFile, os.O_RDONLY, 0644)
		if err != nil {
			totalSize += size
			continue
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
