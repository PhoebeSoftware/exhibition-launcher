package models

import "path/filepath"

type Settings struct {
	DownloadPath   string `json:"download_path"`
	UseRealDebrid bool `json:"use_real_debrid"`
}

func (s Settings) GetSettings() Settings {
	return s
}

func (s *Settings) DefaultValues() {
	s.DownloadPath = filepath.Join("downloads")
	s.UseRealDebrid = false
}