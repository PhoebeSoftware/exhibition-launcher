package main

import (
	"embed"
	"errors"
	"exhibition-launcher/exhibitionQueue"
	"exhibition-launcher/igdb"
	"exhibition-launcher/library"
	"exhibition-launcher/torrent"
	"exhibition-launcher/torrent/realdebrid"
	"exhibition-launcher/utils"
	"exhibition-launcher/utils/jsonUtils"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.
var libraryManager *library.Library
var apiManager *igdb.APIManager
var torrentManager *torrent.Manager
var debridClient *realdebrid.RealDebridClient

//go:embed all:frontend/dist
var assets embed.FS

type WindowService struct{}

var app *application.App

func (w *WindowService) Minimize() {
	win := app.CurrentWindow()
	if win.IsMinimised() {
		win.UnMinimise()
	} else {
		win.Minimise()
	}
}

func (w *WindowService) Maximize() {
	win := app.CurrentWindow()
	if win.IsMaximised() {
		win.UnMaximise()
	} else {
		win.Maximise()
	}
}

func (w *WindowService) Close() {
	app.CurrentWindow().Close()
}

var (
	ErrorTokenIsEmpty = errors.New("Real-debrid token is empty")
)

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// 🐐routine
	settings := &jsonModels.Settings{}
	settingsManager, err := jsonUtils.NewJsonManager(filepath.Join("settings.json"), settings)
	if err != nil {
		fmt.Println(err)
		return
	}

	apiManager, err = igdb.NewAPI(settings, settingsManager)
	if err != nil {
		fmt.Println("error fetching acces token:", err)
		return
	}


	libraryManager = library.GetLibrary(apiManager)
	if settings.RealDebridSettings.UseRealDebrid {
		if settings.RealDebridSettings.DebridToken == "" {
			// TO:DO ADD A UI FOR THIS OR SMTH
			fmt.Println(ErrorTokenIsEmpty)
			return
		}
		debridClient = realdebrid.NewRealDebridClient(settings)
	}
	torrentManager = torrent.StartClient(settings.DownloadPath)
	queue := exhibitionQueue.Queue{
		DownloadsInQueue: []exhibitionQueue.Download{},
		TorrentManager:   torrentManager,
		RealDebridClient: debridClient,
		DownloadPath:     settings.DownloadPath,
		Paused:           false,
	}

	// This code is for refetching covers and banners but it will slow down startup
	//for id, game := range libraryManager.Games {
	//	gameData, err := apiManager.GetGameData(game.IGDBID)
	//	if err != nil {
	//		fmt.Println("Error fetching data for game:", game.IGDBID, err)
	//		continue
	//	}
	//
	//	game.CoverURL = gameData.CoverURL
	//	game.ArtworkUrlList = gameData.ArtworkUrlList
	//	game.ScreenshotUrlList = gameData.ScreenshotUrlList
	//	libraryManager.Games[id] = game
	//}



	webViewWindowOpt := application.WebviewWindowOptions{
		Title:     "Exhibition Launcher",
		Width:     1200,
		Height:    900,
		MinHeight: 700,
		MinWidth:  1064,
		Frameless: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	}

	services := []application.Service{
		//application.NewService(torrentManager),
		application.NewService(apiManager),
		application.NewService(libraryManager),
		application.NewService(&WindowService{}),
		application.NewService(settings),
		application.NewService(settingsManager),
		application.NewService(&utils.PathUtil{}),
		application.NewService(&queue),
	}

	if debridClient != nil {
		services = append(services, application.NewService(debridClient))
	}

	appOptions := application.Options{
		Name:     "Exhibition Launcher",
		Services: services,
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	}
	// If macos turn off frameless so minimizing works
	if runtime.GOOS == "darwin" {
		webViewWindowOpt.Frameless = false

		webViewWindowOpt.MinimiseButtonState = application.ButtonHidden
		webViewWindowOpt.MaximiseButtonState = application.ButtonHidden
		webViewWindowOpt.CloseButtonState = application.ButtonHidden
	}

	app = application.New(appOptions)
	app.NewWebviewWindowWithOptions(webViewWindowOpt)

	queue.App = app
	queue.AddTorrentDownloadToQueue("magnet:?xt=urn:btih:AF8061222ACBB5EDD6DA353E4177FFC403426A88&dn=Schedule+I+v0.3.3f14&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.moeking.me%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.theoks.net%3A6969%2Fannounce&tr=udp%3A%2F%2Fmovies.zsw.ca%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker-udp.gbitt.info%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.gbitt.info%3A80%2Fannounce&tr=https%3A%2F%2Ftracker.gbitt.info%3A443%2Fannounce&tr=http%3A%2F%2Ftracker.ccp.ovh%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.ccp.ovh%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce")


	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
