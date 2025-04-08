package exhibitionQueue

import (
	"exhibition-launcher/utils"
	"fmt"
	"path/filepath"
	"slices"

	"golift.io/xtractr"
)

func ExtractFiles(paths []string, downloadPath string, torrentName string) {
	for _, path := range paths {
		ext := filepath.Ext(path)
		if !slices.Contains(Extensions, ext) {
			continue
		}

		fmt.Println("Extracting file:", path)
		size, files, _, err := xtractr.ExtractFile(&xtractr.XFile{
			FilePath:  filepath.Join(downloadPath, path),
			OutputDir: filepath.Join(downloadPath, torrentName),
		})

		if err != nil {
			fmt.Println("Error extracting file:", err)
			continue
		}

		fmt.Printf("Extracted %d files, final disk usage: %s\n", len(files), utils.HumanizeBytes(float64(size)))
	}
}
