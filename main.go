package main

import (
	"derpy-launcher072/igdb"
	"derpy-launcher072/library"
	"derpy-launcher072/torrent"
	"derpy-launcher072/torrent/realdebrid"
	"derpy-launcher072/utils/jsonUtils"
	"derpy-launcher072/utils/jsonUtils/jsonModels"
	"embed"
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

	libraryManager = library.GetLibrary(apiManager)
	apiManager = igdb.NewAPI()
	if settings.UseRealDebrid {
		if settings.DebridToken == "" {
			// TO:DO ADD A UI FOR THIS OR SMTH
			fmt.Println("Debrid does not exist")
			return
		}
		debridClient = realdebrid.NewRealDebridClient(settings.DebridToken)
	}

	torrentManager = torrent.StartClient(settings.DownloadPath)

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
		Title:     "derpyLauncher",
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
		application.NewService(apiManager),
		application.NewService(libraryManager),
		application.NewService(&WindowService{}),
		application.NewService(settings),
		application.NewService(settingsManager),
	}

	if debridClient != nil {
		services = append(services, application.NewService(debridClient))
	}

	appOptions := application.Options{
		Name: "derp-launcher072",
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

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
