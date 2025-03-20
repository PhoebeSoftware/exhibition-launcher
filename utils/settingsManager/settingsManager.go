package settingsManager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Manager struct {
	Settings Settings
}

type Settings struct {
	PathToSettings string `json:"path_to_settings"`
	DownloadPath   string `json:"download_path"`
	UseRealDebrid bool `json:"use_real_debrid"`
	DebridToken string `json:"debrid_token"`
}

func (settingsManager Manager) GetSettings() Settings {
	return settingsManager.Settings
}

func (settingsManager Manager) SaveSettings() error {
	// Idk random perms idk
	file, err := os.OpenFile(settingsManager.GetSettings().PathToSettings, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not save settingsManager: %w", err)
	}

	jsonData, err := json.MarshalIndent(settingsManager, "", "    ")
	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("could not write json data to settingsManager: %w", err)
	}

	return nil
}

func LoadSettings(path string) (*Manager, error) {
	settingsManager := &Manager{}
	settingsManager.Settings.PathToSettings = path

	file, err := os.Open(path)
	if err != nil {
		settingsManager.Settings.DownloadPath = filepath.Join("downloads")
		err := settingsManager.GenerateSettings()
		if err != nil {
			return nil, err
		}
		file.Close()

		// Reopen file when done generating settingsManager
		file, err = os.Open(path)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(settingsManager)
	if err != nil {
		return nil, fmt.Errorf("Could not decode %v: %w", settingsManager.Settings.PathToSettings, err.Error())
	}
	return settingsManager, nil
}

func (settingsManager *Manager) GenerateSettings() error {
	settingsFile, err := os.Create(settingsManager.Settings.PathToSettings)
	if err != nil {
		return err
	}
	defer settingsFile.Close()

	err = settingsManager.SaveSettings()
	if err != nil {
		return err
	}
	
	settingsManager, err = LoadSettings(settingsManager.Settings.PathToSettings)
	if err != nil {
		return fmt.Errorf("could not load settingsManager: %w", err)
	}
	return nil
}
