package library

import (
	"exhibition-launcher/proxy_client"
	"exhibition-launcher/search"
	"exhibition-launcher/utils/json_utils"
	"exhibition-launcher/utils/json_utils/json_models"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"time"
)

type LibraryManager struct {
	JsonManager *json_utils.JsonManager
	Library     *json_models.Library
	Client      *http.Client
	Settings    *json_models.Settings
	ProxyClient *proxy_client.ProxyClient
	FuzzyManager *search.FuzzyManager
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
func GetLibrary(proxyClient *proxy_client.ProxyClient, settings *json_models.Settings, fuzzyManager *search.FuzzyManager) (*LibraryManager, error) {
	library := &json_models.Library{}
	jsonManager, err := json_utils.NewJsonManager(filepath.Join("library.json"), library)
	if err != nil {
		return nil, err
	}

	fuzzyManager.IndexFuzzy(library.Games)

	return &LibraryManager{
		JsonManager: jsonManager,
		Library:     library,
		ProxyClient: proxyClient,
		Client:      &http.Client{},
		Settings:    settings,
		FuzzyManager: fuzzyManager,
	}, nil
}

func (l *LibraryManager) GetAllGames() map[int]json_models.Game {
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
func (l *LibraryManager) GetGame(igdbId int) (json_models.Game, error) {
	game, ok := l.Library.Games[igdbId]
	if !ok {
		return game, fmt.Errorf("game with IGDB ID %d not found", igdbId)
	}
	return game, nil
}

func (l *LibraryManager) GetRangeGame(amount int, offset int) ([]json_models.Game, error) {
	var games []json_models.Game

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

var isAddingGame bool
func (l *LibraryManager) AddToLibrary(igdbId int) (json_models.Game, error) {
	// prompt executable location
	var (
		game       json_models.Game
		executable = ""
		err        error
	)
	for {
		if !isAddingGame {
			break
		}
		time.Sleep(1 * time.Second)
	}
	isAddingGame = true
	defer func() {
		isAddingGame = false
	}()
	game, ok := l.Library.Games[igdbId]
	if ok && game.Executable != "" {
		fmt.Println("Game is already in library (and does not require path):", game.Name)
		return game, fmt.Errorf("Game already exists")
	}

	/*	dialog := application.OpenFileDialog()
		dialog.SetTitle("Select game executable")
		dialog.AddFilter("Executable files", "*.exe; *.app; *.ink; *.bat;")

		path, err := dialog.PromptForSingleSelection()
		if err != nil {
			fmt.Println("Error selecting file:", err)
			return game, err
		}
		if path == "" {
			fmt.Println("No file selected")
			return game, err
		}

		fmt.Println("Selected file:", path)

	*/
	executable = ""

	// game data
	gameData, err := l.ProxyClient.GetMetadataByIGDBID(igdbId)
	if err != nil {
		return game, err
	}
	if gameData.Name == "" {
		return game, fmt.Errorf("failed to get game data")
	}

	// Append the new game
	game = json_models.Game{
		IGDBID:            igdbId,
		Name:              gameData.Name,
		Description:       gameData.Description,
		PlayTime:          0,
		Achievements:      []int{},
		Executable:        executable,
		Running:           false,
		Favorite:          false,
		CoverURL:          gameData.GetCoverURL(),
		ScreenshotUrlList: gameData.GetScreenshotURLS(),
		ArtworkUrlList:    gameData.GetArtworkURLS(),
	}
	if l.Settings.CacheImagesToDisk {
		// If caching fails still add game to library just with https images instead
		err = l.CacheAllImages(&game)
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
	l.FuzzyManager.IndexFuzzy(l.Library.Games)

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
