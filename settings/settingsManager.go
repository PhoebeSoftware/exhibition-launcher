package settings

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Settings struct {
	PathToSettings string `json:"path_to_settings"`
	DownloadPath   string `json:"download_path"`
}

func LoadSettings(path string) (*Settings, error) {
	settings := &Settings{}
	settings.PathToSettings = path

	file, err := os.Open(path)
	if err != nil {
		settings.DownloadPath = filepath.Join("downloads")
		err := settings.GenerateSettings()
		if err != nil {
			return nil, err
		}
		file.Close()

		// Reopen file when done generating settings
		file, err = os.Open(path)
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(settings)
	if err != nil {
		return nil, fmt.Errorf("Could not decode %v: %w", settings.PathToSettings, err.Error())
	}
	return settings, nil
}

func (settings *Settings) GenerateSettings() error {
	settingsFile, err := os.Create(settings.PathToSettings)
	if err != nil {
		return err
	}
	defer settingsFile.Close()

	jsonData, err := json.MarshalIndent(settings, "", "    ")
	if _, err := settingsFile.Write(jsonData); err != nil {
		return fmt.Errorf("could not write json data to settings: %w", err)
	}

	settings, err = LoadSettings(settings.PathToSettings)
	if err != nil {
		return fmt.Errorf("could not load settings: %w", err)
	}
	return nil
}
