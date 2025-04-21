package main

import (
	"embed"
	"errors"
	"exhibition-launcher/exhibition_queue"
	"exhibition-launcher/library"
	"exhibition-launcher/providers"
	"exhibition-launcher/proxy_client"
	"exhibition-launcher/search"
	"exhibition-launcher/torrent"
	"exhibition-launcher/torrent/realdebrid"
	"exhibition-launcher/utils"
	"exhibition-launcher/utils/json_utils"
	"exhibition-launcher/utils/json_utils/json_models"
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
var libraryManager *library.LibraryManager
var torrentManager *torrent.Manager
var debridClient *realdebrid.RealDebridClient
var providerManager *providers.ProviderManager

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
	settings := &json_models.Settings{}
	_, err := json_utils.NewJsonManager(filepath.Join("settings.json"), settings)
	if err != nil {
		fmt.Println(err)
		return
	}
	proxyClient := proxy_client.NewProxyClient(settings)
	fuzzyManager := &search.FuzzyManager{}
	fmt.Println("The server url is:", proxyClient.BaseURL)
	libraryManager, err = library.GetLibrary(proxyClient, settings, fuzzyManager)
	if err != nil {
		fmt.Println(err)
		return
	}

	if settings.RealDebridSettings.UseRealDebrid {
		if settings.RealDebridSettings.DebridToken == "" {
			// TO:DO ADD A UI FOR THIS OR SMTH
			fmt.Println(ErrorTokenIsEmpty)
			return
		}
		debridClient = realdebrid.NewRealDebridClient(settings)
	}

	torrentManager, err = torrent.StartClient(settings.DownloadPath, settings.BitTorrentSettings)
	if err != nil {
		fmt.Println(err)
		return
	}

	queue := exhibition_queue.Queue{
		DownloadsInQueue: []exhibition_queue.Download{},
		TorrentManager:   torrentManager,
		RealDebridClient: debridClient,
		DownloadPath:     settings.DownloadPath,
		Paused:           false,
		QueueStatus:      exhibition_queue.Idle,
	}


	webViewWindowOpt := application.WebviewWindowOptions{
		Title:     "Exhibition",
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
		application.NewService(torrentManager),
		application.NewService(proxyClient),
		application.NewService(libraryManager),
		application.NewService(&WindowService{}),
		application.NewService(settings),
		application.NewService(&utils.PathUtil{}),
		application.NewService(&queue),
		application.NewService(fuzzyManager),
	}

	if debridClient != nil {
		services = append(services, application.NewService(debridClient))
	}

	appOptions := application.Options{
		Name:     "Exhibition",
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

	var consent_jij = false

	// provider goroutine want het yield
	go func() {
		if !consent_jij {
			return
		}

		providerManager = providers.NewProviderManager()

		// download die sources wrm niet
		for _, sourceLink := range settings.DownloadSources {
			err := providerManager.DownloadProvider(sourceLink)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		// test query
		results := providerManager.SearchDownloadsByGameName("DARK SOULS REMASTERED")
		once := false

		for provider, download := range results {
			fmt.Printf("Provider: %s, Magnets: %d\n", provider, len(download.Magnets))

			if !once {
				queue.AddTorrentDownloadToQueue(download.Magnets[0])
				once = true
			}
		}
		fmt.Printf("results from %d providers found\n", len(results))
	}()


	// Demo games for presentations idk
	/*	go func() {
		listOfCoolGames := []int{
			1905, 1942, 2155, 2368, 7194, 7331, 7334,
			7346, 7360, 7360, 9927, 11133, 11208, 11737, 12517,
			14593, 17000, 19560, 25076, 26192, 26472, 45181,
			75235, 76882, 112875, 113112, 114795, 119133, 119171,
			119277, 119388, 125174, 126098, 135243, 200551, 283363,
		}

		for _, id := range listOfCoolGames {
			game, err := libraryManager.AddToLibrary(id, false)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("Added new game %v : %v\n",id, game.Name)
		}

	}()*/

	if settings.CacheImagesToDisk {
		library.StartImageServer()
		go libraryManager.CheckForCache()
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}

	systray := app.NewSystemTray()
	systray.SetLabel("Exhibition")
	systray.Run()
}
