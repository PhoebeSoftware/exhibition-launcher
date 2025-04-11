package library

import (
	"exhibition-launcher/igdb"
	"exhibition-launcher/utils/jsonUtils"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type LibraryManager struct {
	JsonManager *jsonUtils.JsonManager
	Library     *jsonModels.Library
	APIManager  *igdb.APIManager
	Client      *http.Client
	Settings    *jsonModels.Settings
}

func (l *LibraryManager) GetSortedIDs() []int {
	keys := make([]int, 0, len(l.Library.Games))

	for k := range l.Library.Games {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// geeft library.json als LibraryManager struct vol met data
func GetLibrary(apiManager *igdb.APIManager, settings *jsonModels.Settings) (*LibraryManager, error) {
	library := &jsonModels.Library{}
	jsonManager, err := jsonUtils.NewJsonManager(filepath.Join("library.json"), library)
	if err != nil {
		return nil, err
	}

	return &LibraryManager{
		JsonManager: jsonManager,
		Library:     library,
		APIManager:  apiManager,
		Client:      &http.Client{},
		Settings:    settings,
	}, nil
}

func (l *LibraryManager) GetAllGames() map[int]jsonModels.Game {
	return l.Library.Games
}

func (l *LibraryManager) GetAmountOfGames() int {
	return len(l.Library.Games)
}

func (l *LibraryManager) GetAllGameIDs() []int {
	var intList []int
	for i := range l.Library.Games {
		intList = append(intList, i)
	}
	return intList
}
func (l *LibraryManager) GetGame(igdbId int) (jsonModels.Game, error) {
	game, ok := l.Library.Games[igdbId]
	if !ok {
		return game, fmt.Errorf("game with IGDB ID %d not found", igdbId)
	}
	return game, nil
}

func (l *LibraryManager) GetRangeGame(amount int, offset int) ([]jsonModels.Game, error) {
	var games []jsonModels.Game

	if len(l.Library.Games) == 0 {
		return games, fmt.Errorf("no games in library")
	}

	if offset < 0 || offset >= len(l.Library.Games) {
		return games, fmt.Errorf("offset out of range")
	}

	var keys []int
	for k := range l.Library.Games {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	end := min(offset+amount, len(keys))

	for _, key := range keys[offset:end] {
		games = append(games, l.Library.Games[key])
	}

	return games, nil
}

func (l *LibraryManager) AddToLibrary(igdbId int, promptDialog bool) (jsonModels.Game, error) {
	// prompt executable location
	var (
		game       jsonModels.Game
		executable = ""
		err        error
	)

	game, ok := l.Library.Games[igdbId]
	if ok {
		fmt.Println("Game is already in library:", game.Name)
		return game, fmt.Errorf("Game already exists")
	}

	if promptDialog {
		dialog := application.OpenFileDialog()
		dialog.SetTitle("Select game executable")
		dialog.AddFilter("Executable files", "*.exe; *.app; *.ink; *.bat;")
		path, err := dialog.PromptForSingleSelection()

		if err != nil {
			return game, fmt.Errorf("failed to select executable: %w", err)
		}

		executable = path
	}

	// game data
	gameData, err := l.APIManager.GetGameData(igdbId)
	if err != nil {
		return game, err
	}
	if gameData.Name == "" {
		return game, fmt.Errorf("failed to get game data")
	}

	// Append the new game
	game = jsonModels.Game{
		IGDBID:            igdbId,
		Name:              gameData.Name,
		Description:       gameData.Description,
		PlayTime:          0,
		Achievements:      []int{},
		Executable:        executable,
		Running:           false,
		Favorite:          false,
		CoverURL:          gameData.CoverURL,
		ScreenshotUrlList: gameData.ScreenshotUrlList,
		ArtworkUrlList:    gameData.ArtworkUrlList,
	}
	if l.Settings.UseCaching {
		// If caching fails still add game to library just with https instead
		err = l.CacheAllImages(&game, gameData)
		if err != nil {
			l.Library.Games[igdbId] = game
			err := l.JsonManager.Save()
			if err != nil {
				return game, err
			}
			return game, err
		}
	}

	l.Library.Games[igdbId] = game

	saveErr := l.JsonManager.Save()
	if saveErr != nil {
		return game, fmt.Errorf("failed to save library: %w", saveErr)
	}

	return game, nil
}

func (l *LibraryManager) StartApp(igdbId int) error {
	game := l.Library.Games[igdbId]

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
	l.Library.Games[igdbId] = game
	_ = l.JsonManager.Save()

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
				l.Library.Games[igdbId] = game
				_ = l.JsonManager.Save()

				return
			}
		}
	}()

	return nil
}
