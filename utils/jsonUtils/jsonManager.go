package jsonUtils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configurable interface {
	DefaultValues()
}

type JsonManager struct {
	Config Configurable
	Path   string
}

func NewJsonManager(path string, config Configurable) (*JsonManager, error) {
	manager := &JsonManager{Config: config, Path: path}

	if err := manager.Load(); err != nil {
		manager.Config.DefaultValues()
		if err := manager.Save(); err != nil {
			return nil, err
		}
	}

	return manager, nil
}

func (manager *JsonManager) Load() error {
	file, err := os.Open(manager.Path)
	if err != nil {
		return fmt.Errorf("could not open json file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(manager.Config); err != nil {
		return fmt.Errorf("could not decode JSON: %w", err)
	}

	return nil
}

func (manager *JsonManager) Save() error {
	file, err := os.Create(manager.Path)
	if err != nil {
		return fmt.Errorf("could not save jsonUtils: %w", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(manager.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %w", err)
	}

	_, err = file.Write(jsonData)
	return err
}
