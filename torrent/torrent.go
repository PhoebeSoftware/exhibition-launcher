package torrent

import (
	"errors"
	"exhibition-launcher/utils"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
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
func StartClient(path string, bittorentSettings jsonModels.BitTorrentSettings) (*Manager, error) {
	dirErr := os.MkdirAll(path, os.ModePerm)
	if dirErr != nil {
		return nil, dirErr
	}

	conf := torrent.DefaultConfig
	conf.DataDir = path
	conf.DataDirIncludesTorrentID = false
	conf.Debug = false

	conf.DHTEnabled = bittorentSettings.UseDHT
	conf.PEXEnabled = bittorentSettings.UsePEX

	conf.ResumeOnStartup = false
	conf.Database = filepath.Join(path, "bittorrent.db")

	conf.PortBegin = bittorentSettings.StartPort
	conf.PortEnd = bittorentSettings.EndPort

	session, err := torrent.NewSession(conf)

	fmt.Println("Ports available:", session.Stats().PortsAvailable)

	if err != nil {
		fmt.Println("Error starting torrent client:", err)
		return nil, err
	}

	return &Manager{
		session:     session,
		httpClient:  &http.Client{},
		downloadDir: path,
	}, nil
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
		return t, errors.New("Insufficient disk space")
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
