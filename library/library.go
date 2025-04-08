package library

import (
	"exhibition-launcher/igdb"
	"exhibition-launcher/utils/jsonUtils"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sqweek/dialog"
)

type Library struct {
	LibraryManager *jsonUtils.JsonManager
	Library        *jsonModels.Library
	APIManager     *igdb.APIManager
}

// geeft library.json als Library struct vol met data
func GetLibrary(apiManager *igdb.APIManager) *Library {
	library := &jsonModels.Library{}
	libraryManager, err := jsonUtils.NewJsonManager(filepath.Join("library.json"), library)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &Library{
		LibraryManager: libraryManager,
		Library:        library,
		APIManager:     apiManager,
	}
}

func (lib *Library) GetAllGames() map[int]jsonModels.Game {
	return lib.Library.Games
}

func (lib *Library) GetAllGameIDs() []int{
	var intList []int
	for i := range lib.Library.Games {
		intList = append(intList, i)
	}
	return intList
}
func (lib *Library) GetGame(igdbId int) (jsonModels.Game, error) {
	game, ok := lib.Library.Games[igdbId]
	if !ok {
		return game, fmt.Errorf("game with IGDB ID %d not found", igdbId)
	}
	return game, nil
}

func (lib *Library) GetRangeGame(amount int, offset int) ([]jsonModels.Game, error) {
	games := []jsonModels.Game{}
	if len(lib.Library.Games) == 0 {
		return games, fmt.Errorf("no games in library")
	}

	if offset < 0 || offset >= len(lib.Library.Games) {
		return games, fmt.Errorf("offset out of range")
	}

	if amount <= 0 || amount > len(lib.Library.Games)-offset {
		return games, fmt.Errorf("amount out of range")
	}

	for i := offset; i < offset+amount; i++ {
		foundGame, ok := lib.Library.Games[i]
		if !ok || len(games) >= amount {
			break
		}
		games = append(games, foundGame)
	}
	return games, nil
}

func (lib *Library) AddToLibrary(igdbId int) (jsonModels.Game, error) {
	// prompt executable location
	var game jsonModels.Game
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
	game = jsonModels.Game{
		IGDBID:            igdbId,
		Name:              gameData.Name,
		Description:       gameData.Description,
		PlayTime:          0,
		Achievments:       []int{},
		Executable:        executable,
		Running:           false,
		Favorite:          false,
		CoverURL:          gameData.CoverURL,
		ScreenshotUrlList: gameData.ScreenshotUrlList,
		ArtworkUrlList:    gameData.ArtworkUrlList,
	}
	lib.Library.Games[igdbId] = game

	saveErr := lib.LibraryManager.Save()
	if saveErr != nil {
		return game, fmt.Errorf("failed to save library: %w", saveErr)
	}

	return game, nil
}

func (lib *Library) StartApp(igdbId int) error {
	game := lib.Library.Games[igdbId]

	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", game.Executable)
	} else {
		cmd = exec.Command(game.Executable)
	}

	cmd.Dir = filepath.Dir(game.Executable)
	err := cmd.Start()
	if err != nil {
		return err
	}

	fmt.Printf("Started game with PID: %d\n", cmd.Process.Pid)

	game.Running = true
	lib.Library.Games[igdbId] = game
	_ = lib.LibraryManager.Save()

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
				lib.Library.Games[igdbId] = game
				_ = lib.LibraryManager.Save()

				return
			}
		}
	}()

	return nil
}
