package torrent

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent"
)

type DownloadData struct {
	Name     string
	Progress int
	Speed    int
}

type Manager struct {
	client *torrent.Client
	games  map[string]DownloadData
}

// start client en geef manager zodat je makkelijk kan bedienen zawg
func StartClient() *Manager {
	torrentPath := filepath.Join(".", "downloads")

	dirErr := os.MkdirAll(torrentPath, os.ModePerm)
	if dirErr != nil {
		fmt.Println("Error creating downloads directory")
	}

	client, err := torrent.NewClient(nil)

	if err != nil {
		fmt.Println("Error starting torrent client")
	}

	return &Manager{client: client, games: make(map[string]DownloadData)}
}

// add download
// start ook torrent meteen
func (manager Manager) AddTorrent(magnetLink string) (*torrent.Torrent, error) {
	t, err := manager.client.AddMagnet(magnetLink)
	if err != nil {
		return t, err
	}

	fmt.Println("Getting metadata")
	<-t.GotInfo()

	fmt.Println("Download starting")
	t.DownloadAll()

	manager.games[t.Info().Name] = DownloadData{
		Name:     t.Info().Name,
		Progress: 0,
		Speed:    0,
	}

	// speed goroutine
	go func() {
		var lastBytes int = 0

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				currentBytes := int(t.BytesCompleted())
				completionRatio := float64(currentBytes) / float64(t.Info().TotalLength())

				game := manager.games[t.Info().Name]
				game.Speed = currentBytes - lastBytes
				game.Progress = int(completionRatio * 100)

				manager.games[t.Info().Name] = game

				lastBytes = currentBytes
				if completionRatio >= 1.0 {
					delete(manager.games, t.Info().Name)
					return
				}
			}
		}
	}()

	return t, nil
}
