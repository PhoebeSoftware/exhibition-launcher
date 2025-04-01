package jsonModels

import "path/filepath"

type Settings struct {
	DownloadPath       string             `json:"download_path"`
	RealDebridSettings RealDebridSettings `json:"real_debrid_settings"`
}

type RealDebridSettings struct {
	UseRealDebrid bool   `json:"use_real_debrid"`
	DebridToken   string `json:"debrid_token"`
}

func (s Settings) GetSettings() Settings {
	return s
}

func (s *Settings) DefaultValues() {
	s.DownloadPath = filepath.Join("downloads")
	s.RealDebridSettings.UseRealDebrid = false
	s.RealDebridSettings.DebridToken = ""
}
