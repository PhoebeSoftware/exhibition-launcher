package torrent

import (
	"errors"
	"exhibition-launcher/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/rain/torrent"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Manager struct {
	session     *torrent.Session
	httpClient  *http.Client
	downloadDir string
}

// start client en geef manager zodat je makkelijk kan bedienen zawg
func StartClient(path string, pex bool, dht bool, port uint16) *Manager {
	dirErr := os.MkdirAll(path, os.ModePerm)
	if dirErr != nil {
		fmt.Println("Error creating downloads directory:", dirErr)
		return nil
	}

	conf := torrent.DefaultConfig
	conf.DataDir = path
	conf.DataDirIncludesTorrentID = false
	conf.Debug = false

	conf.DHTEnabled = dht
	conf.PEXEnabled = pex

	conf.ResumeOnStartup = false
	conf.Database = filepath.Join(path, "bittorrent.db")

	// alleen frigging ranges zijn er
	conf.PortBegin = port
	conf.PortEnd = port + 10

	session, err := torrent.NewSession(conf)

	fmt.Println("port available:", session.Stats().PortsAvailable)

	if err != nil {
		fmt.Println("Error starting torrent client:", err)
		return nil
	}

	return &Manager{
		session:     session,
		httpClient:  &http.Client{},
		downloadDir: path,
	}
}

func (manager *Manager) RemoveTorrent(uuid string) error {
	err := manager.session.RemoveTorrent(uuid)
	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) SetPaused(value bool) {
	fmt.Println("Setting BitTorrent paused to", value)

	if value {
		manager.session.StopAll()
	} else {
		manager.session.StartAll()
	}
}

func IsTriggered(ch <-chan struct{}) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

// add download
// start ook torrent meteen
func (manager Manager) AddTorrent(app *application.App, uuid string, magnetLink string) (*torrent.Torrent, error) {
	startTime := time.Now()

	// add torrent to session
	fmt.Println("Adding torrent", uuid)
	t, err := manager.session.AddURI(magnetLink, &torrent.AddTorrentOptions{
		ID: uuid,
	})
	if err != nil {
		fmt.Println("Error adding torrent:", err)
		return t, err
	}

	fmt.Println("using port:", t.Port())

	// get metadata
	for !IsTriggered(t.NotifyMetadata()) {
		fmt.Println("Searching for metadata...")
		time.Sleep(time.Second)
	}

	// get space data
	disk := utils.DiskUsage(manager.downloadDir)
	sizeLeft := t.Stats().Bytes.Incomplete

	fmt.Printf("%s free (%s left to download)\n", utils.HumanizeBytes(float64(disk.Free)), utils.HumanizeBytes(float64(sizeLeft)))

	// check space
	if int64(disk.Free) < sizeLeft {
		manager.session.RemoveTorrent(uuid)
		fmt.Println("Insufficient disk space")
		return t, errors.New("get yo bread right nigga")
	}

	fmt.Println("BitTorrent download starting")

	// speed goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		lastBytesDown := float64(0)
		lastBytesUp := float64(0)

		for range ticker.C {
			stats := t.Stats()
			if stats.Status == torrent.Stopped || stats.Status == torrent.Seeding {
				if IsTriggered(t.NotifyComplete()) || IsTriggered(t.NotifyClose()) {
					fmt.Printf("BitTorrent download complete for %s!\n", t.Name())
					app.EmitEvent("download_complete", "Download Finished!")
					break
				}
				fmt.Println("Torrent paused")

				continue
			}

			if stats.Status != torrent.Stopping {
				completionRatio := float64(stats.Bytes.Completed) / float64(stats.Bytes.Total)
				percentage := int(completionRatio * 100)

				downSpeed := utils.HumanizeBytes(float64(stats.Bytes.Completed) - float64(lastBytesDown))
				upSpeed := utils.HumanizeBytes(float64(stats.Bytes.Uploaded) - float64(lastBytesUp))

				lastBytesDown = float64(stats.Bytes.Completed)
				lastBytesUp = float64(stats.Bytes.Uploaded)

				fmt.Printf("\n%s\nProgress: %d%%\n↑: %s/s\n↓: %s/s\nETA: %s\n", t.Name(), percentage, upSpeed, downSpeed, stats.ETA)
				app.EmitEvent("download_progress", map[string]interface{}{
					"percent":         percentage,
					"downloadedBytes": stats.Bytes.Completed,
					"totalBytes":      stats.Bytes.Total,
					"timePassed":      time.Since(startTime).String(),
				})
			}
		}
	}()

	return t, nil
}
