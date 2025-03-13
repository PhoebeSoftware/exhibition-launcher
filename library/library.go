package library

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Game struct {
	AppID       int    `json:"appid"`
	PlayTime    int    `json:"playtime"`
	Achievments []int  `json:"achievments"`
	Executable  string `json:"executable"`
	Running     bool   `json:"running"`
	Favorite    bool   `json:"favorite"`
}

type Library struct {
	Games map[int]Game `json:"games"`
}

// geeft library.json als Library struct vol met data
func GetLibrary() *Library {
	file, err := os.OpenFile(filepath.Join(".", "downloads"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error opening/creating library.json: %v", err)
		return &Library{Games: make(map[int]Game)}
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading library.json: %v", err)
		return &Library{Games: make(map[int]Game)}
	}

	if len(bytes) == 0 {
		emptyLib := &Library{Games: make(map[int]Game)}
		jsonData, err := json.MarshalIndent(emptyLib, "", "    ")
		if err != nil {
			log.Printf("Error marshaling empty library: %v", err)
			return emptyLib
		}
		if _, err := file.Write(jsonData); err != nil {
			log.Printf("Error writing empty library: %v", err)
		}
		return emptyLib
	}

	var library Library
	if err := json.Unmarshal(bytes, &library); err != nil {
		log.Printf("Error unmarshalling library.json: %v", err)
		return &Library{Games: make(map[int]Game)}
	}

	return &library
}

func (lib *Library) AddToLibrary(gameData Game) error {
	// Append the new game
	lib.Games[gameData.AppID] = gameData

	// Marshal the entire library to JSON
	jsonData, err := json.Marshal(lib)
	if err != nil {
		return fmt.Errorf("failed to marshal library: %w", err)
	}

	// Write to file
	err = os.WriteFile("library.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write library file: %w", err)
	}

	return nil
}

func (lib *Library) StartApp(appID int) bool {
	game := lib.Games[appID]

	cmd := exec.Command(game.Executable)
	cmd.Dir = filepath.Dir(game.Executable)
	cmd.Start()

	fmt.Printf("Started game with PID: %d\n", cmd.Process.Pid)

	game.Running = true
	lib.AddToLibrary(game)

	go func() {
		seconds := 0
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("Game running for %d seconds\n", seconds)
				seconds++
			case <-done:
				fmt.Printf("Game quit after %d seconds\n", seconds)

				game.Running = false
				game.PlayTime += seconds
				lib.AddToLibrary(game)
				return
			}
		}
	}()

	return true
}
