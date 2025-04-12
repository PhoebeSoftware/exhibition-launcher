package main

import (
	"embed"
	"errors"
	"exhibition-launcher/exhibitionQueue"
	"exhibition-launcher/igdb"
	"exhibition-launcher/library"
	"exhibition-launcher/providers"
	"exhibition-launcher/search"
	"exhibition-launcher/torrent"
	"exhibition-launcher/torrent/realdebrid"
	"exhibition-launcher/utils"
	"exhibition-launcher/utils/jsonUtils"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.
var libraryManager *library.LibraryManager
var igdbApiManager *igdb.APIManager
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
	settings := &jsonModels.Settings{}
	settingsManager, err := jsonUtils.NewJsonManager(filepath.Join("settings.json"), settings)
	if err != nil {
		fmt.Println(err)
		return
	}
	// test code
	settings.UseDirectIGDB = true
	if settings.UseDirectIGDB {
		igdbApiManager, err = igdb.NewAPI(settings, settingsManager)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	libraryManager, err = library.GetLibrary(igdbApiManager, settings)
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

	queue := exhibitionQueue.Queue{
		DownloadsInQueue: []exhibitionQueue.Download{},
		TorrentManager:   torrentManager,
		RealDebridClient: debridClient,
		DownloadPath:     settings.DownloadPath,
		Paused:           false,
		QueueStatus:      exhibitionQueue.Idle,
	}

	searchManager := search.SearchManager{
		LibraryManager: libraryManager,
	}

	// Always index 
	searchManager.IndexGames()
	// Demo
	fmt.Println("Total games:", len(libraryManager.Library.Games))
	startTime := time.Now()
	ids := searchManager.SearchForName("bloodborn")
	for _, id := range ids {
		game := libraryManager.Library.Games[id]
		fmt.Printf("%v:%v\n", game.Name, id)
	}
	fmt.Println(time.Since(startTime))
	

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
		application.NewService(torrentManager),
		application.NewService(igdbApiManager),
		application.NewService(libraryManager),
		application.NewService(&WindowService{}),
		application.NewService(settings),
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

	var consent_jij = false

	// provider goroutine want het yield
	go func() {
		if !consent_jij {
			return
		}

		providerManager = providers.NewProviderManager()

		// download die sources wrm niet
		for _, sourceLink := range settings.DownloadSources {
			providerManager.DownloadProvider(sourceLink)
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

//	Add a bunch of games

	//go func() {
	//	for i := 2000; i < 7000; i++ {
	//		game, err := libraryManager.AddToLibrary(i, false)
	//		if err != nil {
	//			fmt.Println(err)
	//			continue
	//		}
	//		fmt.Println("Added game:", game.Name)
	//	}
	//}()

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

	if settings.UseCaching {
		library.StartImageServer()
		go libraryManager.CheckForCache()
	}

	// Complete aot season 4
	queue.AddRealDebridDownloadToQueue("magnet:?xt=urn:btih:ac8dc037d282f82efb2864abdd54399029105c0c&dn=%5BGolumpa-Yameii%5D%20Attack%20on%20Titan%20-%20The%20Final%20Season%20%5BEnglish%20Dub%5D%20%5BWEB-DL%20720p%5D%20-%20%28The%20Complete%20S04%29%20-%20Unofficial%20Batch&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce")

	//  Various Artists - NOW – Yearbook The Vault 1980 (2025) Mp3 320kbp...
	queue.AddRealDebridDownloadToQueue("magnet:?xt=urn:btih:BE08A4F593706D533D7A8B7BEEB5477EBD2D3F4F&dn=Various+Artists+-+NOW+%26ndash%3B+Yearbook+The+Vault+1980+%282025%29+Mp3+320kbps+%5BPMEDIA%5D+%E2%AD%90%EF%B8%8F&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fopentracker.io%3A6969%2Fannounce&tr=udp%3A%2F%2Fz.mercax.com%3A53%2Fannounce&tr=udp%3A%2F%2Ftracker.birkenwald.de%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.bittor.pw%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.renfei.net%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce")

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
