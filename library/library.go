package library

import (
	"encoding/json"
	"exhibition-launcher/igdb"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sqweek/dialog"
)

type Game struct {
	IGDBID      int    `json:"igdb_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlayTime    int    `json:"playtime"`
	Achievments []int  `json:"achievments"`
	Executable  string `json:"executable"`
	Running     bool   `json:"running"`
	Favorite    bool   `json:"favorite"`
	MainCover   string
	Banners      []string
}

type Library struct {
	Games      map[int]Game `json:"games"`
	APIManager *igdb.APIManager
}

// geeft library.json als Library struct vol met data
func GetLibrary(apiManager *igdb.APIManager) *Library {
	file, err := os.OpenFile(filepath.Join(".", "library.json"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error opening/creating library.json: %v", err)
		return &Library{
			Games:      make(map[int]Game),
			APIManager: apiManager,
		}
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading library.json: %v", err)
		return &Library{
			Games:      make(map[int]Game),
			APIManager: apiManager,
		}
	}

	if len(bytes) == 0 {
		emptyLib := &Library{
			Games:      make(map[int]Game),
			APIManager: apiManager,
		}
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
		return &Library{
			Games:      make(map[int]Game),
			APIManager: apiManager,
		}
	}

	library.APIManager = apiManager
	return &library
}

func (lib *Library) GetAllGames() map[int]Game {
	return lib.Games
}

func (lib *Library) AddToLibrary(igdbId int) (Game, error) {
	// prompt executable location
	var game Game
	executable, err := dialog.File().Title("Select game executable").Filter("Executable files", "exe", "app", "ink", "bat").Load()
	if err != nil {
		return game, fmt.Errorf("failed to select executable: %w", err)
	}

	// game data
	gameData, err := lib.APIManager.GetGameData(igdbId)
	if err != nil {
		return game, err
	}
	fmt.Println("Game data:")
	fmt.Println(gameData)
	if gameData.Name == "" {
		return game, fmt.Errorf("failed to get game data")
	}

	// Append the new game
	game = Game{
		IGDBID:      igdbId,
		Name:        gameData.Name,
		Description: gameData.Description,
		PlayTime:    0,
		Achievments: []int{},
		Executable:  executable,
		Running:     false,
		Favorite:    false,
		MainCover:   gameData.MainCover,
		Banners: gameData.Banners,
	}
	lib.Games[igdbId] = game

	// Marshal the entire library to JSON
	jsonData, err := json.Marshal(lib)
	if err != nil {
		return game, fmt.Errorf("failed to marshal library: %w", err)
	}

	// Write to file
	err = os.WriteFile("library.json", jsonData, 0644)
	if err != nil {
		return game, fmt.Errorf("failed to write library file: %w", err)
	}

	return game, nil
}

func (lib *Library) StartApp(igdbId int) bool {
	game := lib.Games[igdbId]

	cmd := exec.Command(game.Executable)
	cmd.Dir = filepath.Dir(game.Executable)
	cmd.Start()

	fmt.Printf("Started game with PID: %d\n", cmd.Process.Pid)

	game.Running = true

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
				return
			}
		}
	}()

	return true
}
