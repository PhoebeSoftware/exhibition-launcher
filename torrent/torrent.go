package torrent

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/rain/torrent"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type DownloadData struct {
	Name     string
	Progress int
	Speed    int
}
type Manager struct {
	session    *torrent.Session
	games      map[string]DownloadData
	httpClient *http.Client
}

// start client en geef manager zodat je makkelijk kan bedienen zawg
func StartClient(path string) *Manager {
	dirErr := os.MkdirAll(path, os.ModePerm)
	if dirErr != nil {
		fmt.Println("Error creating downloads directory")
		return nil
	}

	conf := torrent.DefaultConfig
	conf.DataDir = path

	session, err := torrent.NewSession(conf)

	if err != nil {
		fmt.Println("Error starting torrent client:", err)
		return nil
	}

	return &Manager{
		session:    session,
		games:      make(map[string]DownloadData),
		httpClient: &http.Client{},
	}
}

func (manager *Manager) SetPaused(value bool) {
	fmt.Println("Setting paused:", value)

	if value {
		manager.session.StopAll()
	} else {
		manager.session.StartAll()
	}

}

// add download
// start ook torrent meteen
func (manager Manager) AddTorrent(app *application.App, magnetLink string) (*torrent.Torrent, error) {
	startTime := time.Now()

	fmt.Println("Adding torrent:", magnetLink)
	t, err := manager.session.AddURI(magnetLink, nil)
	if err != nil {
		return t, err
	}

	fmt.Println("Getting metadata")
	<-t.NotifyMetadata()
	fmt.Println("Download starting")

	manager.games[t.Name()] = DownloadData{
		Name:     t.Name(),
		Progress: 0,
		Speed:    0,
	}

	// speed goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				stats := t.Stats()

				completionRatio := float64(stats.Bytes.Completed) / float64(stats.Bytes.Total)

				game := manager.games[t.Name()]
				game.Speed = stats.Speed.Download
				game.Progress = int(completionRatio * 100)

				manager.games[t.Name()] = game
				fmt.Printf("Game: %s, Progress: %d%%, Speed: %d bytes/s\n", game.Name, game.Progress, game.Speed)
				app.EmitEvent("download_progress", map[string]interface{}{
					"percent":         game.Progress,
					"downloadedBytes": stats.Bytes.Completed,
					"totalBytes":      stats.Bytes.Total,
					"timePassed":      time.Since(startTime).String(),
				})

				if completionRatio >= 1.0 {
					fmt.Printf("Download complete: %s\n", t.Name())
					app.EmitEvent("download_complete", "Download Finished!")

					delete(manager.games, t.Name())
					return
				}
			}
		}
	}()

	return t, nil
}
