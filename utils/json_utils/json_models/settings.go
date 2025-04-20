package json_models

import (
	"path/filepath"
	"time"
)

type Settings struct {
	DownloadPath       string             `json:"download_path"`
	DownloadSources    []string           `json:"download_sources"`
	UseDirectIGDB      bool               `json:"use_direct_igdb"`
	UseCaching         bool               `json:"use_caching"`
	RealDebridSettings RealDebridSettings `json:"real_debrid_settings"`
	IgdbSettings       IgdbSettings       `json:"igdb_settings"`
	BitTorrentSettings BitTorrentSettings `json:"bittorrent_settings"`
}

type RealDebridSettings struct {
	UseRealDebrid   bool   `json:"use_real_debrid"`
	DebridToken     string `json:"debrid_token"`
	NumberOfThreads int    `json:"number_of_threads"`
}

type IgdbSettings struct {
	IgdbClient string `json:"igdb_client"`
	IgdbSecret string `json:"igdb_secret"`
	IgdbAuth   string `json:"igdb_auth"`

	// In seconds
	ExpiresIn int `json:"expires_in"`

	// Basic go time format
	ExpiresAt time.Time `json:"expires_at"`
}

type BitTorrentSettings struct {
	UseDHT    bool   `json:"use_dht"`
	UsePEX    bool   `json:"use_pex"`
	StartPort uint16 `json:"start_port"`
	EndPort   uint16 `json:"end_port"`
}

func (s Settings) GetSettings() Settings {
	return s
}

func (s *Settings) DefaultValues() {
	s.DownloadPath = filepath.Join("downloads")
	s.UseDirectIGDB = true
	s.UseCaching = false
	s.DownloadSources = []string{}

	s.RealDebridSettings.UseRealDebrid = false
	s.RealDebridSettings.DebridToken = ""
	s.RealDebridSettings.NumberOfThreads = 2

	s.IgdbSettings.IgdbClient = "client_id"
	s.IgdbSettings.IgdbSecret = "client_secret"
	s.IgdbSettings.IgdbAuth = "auto_generated_on_launch"
	s.IgdbSettings.ExpiresIn = 0

	s.BitTorrentSettings.UseDHT = true
	s.BitTorrentSettings.UsePEX = true
	s.BitTorrentSettings.StartPort = 9000
	s.BitTorrentSettings.EndPort = 9010
}
