package exhibition_queue

import (
	"exhibition-launcher/utils"
	"fmt"
	"path/filepath"
	"slices"

	"golift.io/xtractr"
)

var Extensions = []string{
	".zip",
	".rar",
	".7z",
}

func ExtractFile(path string, downloadPath string, torrentName string) error {
	ext := filepath.Ext(path)
	if !slices.Contains(Extensions, ext) {
		return nil
	}

	fmt.Println("Extracting file:", path)
	size, files, _, err := xtractr.ExtractFile(&xtractr.XFile{
		FilePath:  filepath.Join(downloadPath, path),
		OutputDir: filepath.Join(downloadPath, torrentName),
	})

	if err != nil {
		fmt.Println("Error extracting file:", err)
		return err
	}

	fmt.Printf("Extracted %d files, final disk usage: %s\n", len(files), utils.HumanizeBytes(float64(size)))
	return nil
}
