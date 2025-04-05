package exhibitionQueue

import (
	"exhibition-launcher/utils"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/cenkalti/rain/torrent"
	"golift.io/xtractr"
)

func ExtractFiles(files []torrent.File, downloadPath string, torrentName string) {
	for _, file := range files {
		ext := filepath.Ext(file.Path())
		if !slices.Contains(Extensions, ext) {
			continue
		}

		fmt.Println("Extracting file:", file.Path())
		size, files, _, err := xtractr.ExtractFile(&xtractr.XFile{
			FilePath:  filepath.Join(downloadPath, file.Path()),
			OutputDir: filepath.Join(downloadPath, torrentName),
		})

		if err != nil {
			fmt.Println("Error extracting file:", err)
			continue
		}

		fmt.Printf("Extracted %d files, final disk usage: %s\n", len(files), utils.HumanizeBytes(float64(size)))
	}
}
