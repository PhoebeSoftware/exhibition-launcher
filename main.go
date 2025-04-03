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

	apiManager = igdb.NewAPI()
	libraryManager = library.GetLibrary(apiManager)
	if settings.RealDebridSettings.UseRealDebrid {
		if settings.RealDebridSettings.DebridToken == "" {
			// TO:DO ADD A UI FOR THIS OR SMTH
			fmt.Println(ErrorTokenIsEmpty)
			return
		}
		debridClient = realdebrid.NewRealDebridClient(settings)
	}

	queue := exhibitionQueue.Queue{
		DownloadsInQueue: []exhibitionQueue.Download{},
		TorrentManager:   torrentManager,
		RealDebridClient: debridClient,
		DownloadPath:     settings.DownloadPath,
		Paused: false,
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

	//torrentManager = torrent.StartClient(settings.DownloadPath)

	//go func() {
	//	results := torrent.Scrape_1337x("goat simulator 3")
	//	for _, result := range results {
	//		data := torrent.Get_1337x_data(result)
	//		fmt.Printf("Title: %s\nUploader: %s\nDownloads: %d\nDate: %s\n\n", data.Title, data.Uploader, data.Downloads, data.Date)
	//	}
	//}()

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.

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

	//go func() {
	//	err := debridClient.DownloadByMagnet("magnet:?xt=urn:btih:EEEF75F8C7AD79818C54C618E1A7937CD76B59C4&dn=Sony+Vegas+Pro+v11.0.510+64+bit+%28patch+keygen+DI%29+%5BChingLiu%5D&tr=http%3A%2F%2Fpow7.com%2Fannounce&tr=http%3A%2F%2Fpubt.net%3A2710%2Fannounce&tr=http%3A%2F%2Ft1.pow7.com%2Fannounce&tr=http%3A%2F%2Ftracker.torrentbay.to%3A6969%2Fannounce&tr=http%3A%2F%2Ftracker.torrent.to%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%2Fannounce&tr=udp%3A%2F%2Ftracker.1337x.org%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.istole.it%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Fa.tracker.prq.to%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce", settings.DownloadPath)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//}()
	app = application.New(appOptions)
	app.NewWebviewWindowWithOptions(webViewWindowOpt)
	queue.App = app
	// Test code
	//  Various Artists - NOW That’s What I Call Jukebox Classics True Lo...
	queue.AddRealDebridDownloadToQueue("magnet:?xt=urn:btih:0B2BB9C69CA59FCBFDB3BE140A6D189D82BA29C3&dn=Various+Artists+-+NOW+That%26rsquo%3Bs+What+I+Call+Jukebox+Classics+True+Love+Ways+%282025%29+Mp3+320kbps+%5BPMEDIA%5D+%E2%AD%90%EF%B8%8F&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fopentracker.io%3A6969%2Fannounce&tr=udp%3A%2F%2Fz.mercax.com%3A53%2Fannounce&tr=udp%3A%2F%2Ftracker.birkenwald.de%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.bittor.pw%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.renfei.net%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce")
	//  Various Artists - NOW – Yearbook The Vault 1980 (2025) Mp3 320kbp...
	queue.AddRealDebridDownloadToQueue("magnet:?xt=urn:btih:BE08A4F593706D533D7A8B7BEEB5477EBD2D3F4F&dn=Various+Artists+-+NOW+%26ndash%3B+Yearbook+The+Vault+1980+%282025%29+Mp3+320kbps+%5BPMEDIA%5D+%E2%AD%90%EF%B8%8F&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fopentracker.io%3A6969%2Fannounce&tr=udp%3A%2F%2Fz.mercax.com%3A53%2Fannounce&tr=udp%3A%2F%2Ftracker.birkenwald.de%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.bittor.pw%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.renfei.net%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce")
	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
