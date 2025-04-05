package realdebrid

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

// Ez chatgpt
func GetMagnetLinkHash(magnetLink string) (string, error) {
	// Regular expression to extract the hash from the magnet link
	re := regexp.MustCompile(`xt=urn:btih:([a-fA-F0-9]{40})`)

	// Match the regex to the magnet link
	matches := re.FindStringSubmatch(magnetLink)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("invalid magnet link or hash not found")
}
