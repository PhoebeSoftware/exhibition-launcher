package main

import (
	"embed"
	"errors"
	"exhibition-launcher/exhibitionQueue"
	"exhibition-launcher/igdb"
	"exhibition-launcher/library"
	"exhibition-launcher/providers"
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
	if settings.UseCaching {
		go libraryManager.CheckForCache()
	}

	if settings.RealDebridSettings.UseRealDebrid {
		if settings.RealDebridSettings.DebridToken == "" {
			// TO:DO ADD A UI FOR THIS OR SMTH
			fmt.Println(ErrorTokenIsEmpty)
			return
		}
		debridClient = realdebrid.NewRealDebridClient(settings)
	}

	torrentManager = torrent.StartClient(settings.DownloadPath, settings.BitTorrentSettings.UsePEX, settings.BitTorrentSettings.UseDHT, settings.BitTorrentSettings.Port)

	// provider goroutine want het yield
	go func() {
		providerManager = providers.NewProviderManager()

		// cache die sources wrm niet
		for _, sourceLink := range settings.DownloadSources {
			providerManager.CacheProvider(sourceLink)
		}

		// test query
		results := providerManager.SearchDownloadsByGameName("Palworld")
		for _, result := range results {
			fmt.Println(result.Magnets)
		}
	}()

	queue := exhibitionQueue.Queue{
		DownloadsInQueue: []exhibitionQueue.Download{},
		TorrentManager:   torrentManager,
		RealDebridClient: debridClient,
		DownloadPath:     settings.DownloadPath,
		Paused:           false,
		QueueStatus:      exhibitionQueue.Idle,
	}

	// This code is for refetching covers and banners but it will slow down startup
	//for id, game := range libraryManager.Games {
	//	gameData, err := igdbApiManager.GetGameData(game.IGDBID)
	//	if err != nil {
	//		fmt.Println("Error fetching data for game:", game.IGDBID, err)
	//		continue
	//	}
	//
	//	game.CoverFilename = gameData.CoverFilename
	//	game.ArtworkFilenames = gameData.ArtworkFilenames
	//	game.ScreenshotFilenames = gameData.ScreenshotFilenames
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

	// Add a bunch of games

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

	// Schedule I v0.3.3f14
	queue.AddTorrentDownloadToQueue("magnet:?xt=urn:btih:7027B6E2A1E4566B0B0DC4146E8B54235B743AE6&dn=Schedule+I+%28v0.3.3f15+%2B+Online+Multiplayer%29+%5BDODI+Repack%5D&tr=udp%3A%2F%2F9.rarbg.to%3A2870%2Fannounce&tr=udp%3A%2F%2Feddie4.nl%3A6969%2Fannounce&tr=udp%3A%2F%2Fthetracker.org%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337&tr=udp%3A%2F%2F9.rarbg.com%3A2710%2Fannounce&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker1.wasabii.com.tw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker.cypherpunks.ru%3A6969%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce")

	// Einstein's Cats v0.1.0
	queue.AddTorrentDownloadToQueue("magnet:?xt=urn:btih:1A2ED74D4E45E9FBACE0B24571051840D81045A4&dn=Einstein%26%23039%3Bs+Cats+v0.1.0&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.moeking.me%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.theoks.net%3A6969%2Fannounce&tr=udp%3A%2F%2Fmovies.zsw.ca%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker-udp.gbitt.info%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.gbitt.info%3A80%2Fannounce&tr=https%3A%2F%2Ftracker.gbitt.info%3A443%2Fannounce&tr=http%3A%2F%2Ftracker.ccp.ovh%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.ccp.ovh%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce")

	// Complete aot season 4
	queue.AddRealDebridDownloadToQueue("magnet:?xt=urn:btih:ac8dc037d282f82efb2864abdd54399029105c0c&dn=%5BGolumpa-Yameii%5D%20Attack%20on%20Titan%20-%20The%20Final%20Season%20%5BEnglish%20Dub%5D%20%5BWEB-DL%20720p%5D%20-%20%28The%20Complete%20S04%29%20-%20Unofficial%20Batch&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce")

	//  Various Artists - NOW – Yearbook The Vault 1980 (2025) Mp3 320kbp...
	queue.AddRealDebridDownloadToQueue("magnet:?xt=urn:btih:BE08A4F593706D533D7A8B7BEEB5477EBD2D3F4F&dn=Various+Artists+-+NOW+%26ndash%3B+Yearbook+The+Vault+1980+%282025%29+Mp3+320kbps+%5BPMEDIA%5D+%E2%AD%90%EF%B8%8F&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fopentracker.io%3A6969%2Fannounce&tr=udp%3A%2F%2Fz.mercax.com%3A53%2Fannounce&tr=udp%3A%2F%2Ftracker.birkenwald.de%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.bittor.pw%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.renfei.net%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce")

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
