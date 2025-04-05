package torrent

import (
	"errors"
	"exhibition-launcher/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/rain/torrent"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

type Manager struct {
	session     *torrent.Session
	httpClient  *http.Client
	downloadDir string
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
	conf.DataDirIncludesTorrentID = false
	conf.Debug = false

	session, err := torrent.NewSession(conf)

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

func HumanizeBytes(bytesPerSec float64) string {
	switch {
	case bytesPerSec >= TB:
		return fmt.Sprintf("%.2f TB", bytesPerSec/TB)
	case bytesPerSec >= GB:
		return fmt.Sprintf("%.2f GB", bytesPerSec/GB)
	case bytesPerSec >= MB:
		return fmt.Sprintf("%.2f MB", bytesPerSec/MB)
	case bytesPerSec >= KB:
		return fmt.Sprintf("%.2f KB", bytesPerSec/KB)
	default:
		return fmt.Sprintf("%.0f B", bytesPerSec)
	}
}

func (manager *Manager) SetPaused(value bool) {
	fmt.Println("Setting BitTorrent paused to ", value)

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

	fmt.Println("Adding torrent ", uuid)
	t, err := manager.session.AddURI(magnetLink, &torrent.AddTorrentOptions{
		ID:                uuid,
		StopAfterDownload: true, // FUCK seeding
	})
	if err != nil {
		return t, err
	}

	// get metadata
	fmt.Println("Getting BitTorrent metadata")
	<-t.NotifyMetadata()

	// check space
	disk := utils.DiskUsage(manager.downloadDir)
	torrentSize := t.Stats().Bytes.Incomplete

	fmt.Printf("%s free (%s left to download)\n", HumanizeBytes(float64(disk.Free)), HumanizeBytes(float64(torrentSize)))

	if int64(disk.Free) < torrentSize {
		fmt.Println("Insufficient disk space")
		return t, errors.New("get yo bread right nigga")
	}

	fmt.Println("BitTorrent Download starting")

	// speed goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			stats := t.Stats()
			if stats.Status == torrent.Stopped || stats.Status == torrent.Stopping {
				continue // paused
			}

			completionRatio := float64(stats.Bytes.Completed) / float64(stats.Bytes.Total)
			percentage := int(completionRatio * 100)
			speedySpeed := HumanizeBytes(float64(stats.Speed.Download))

			fmt.Printf("%s: Progress: %d%%, Speed: %s/s ETA: %s\n", t.Name(), percentage, speedySpeed, stats.ETA)
			app.EmitEvent("download_progress", map[string]interface{}{
				"percent":         percentage,
				"downloadedBytes": stats.Bytes.Completed,
				"totalBytes":      stats.Bytes.Total,
				"timePassed":      time.Since(startTime).String(),
			})

			// trust
			if IsTriggered(t.NotifyComplete()) || IsTriggered(t.NotifyClose()) {
				fmt.Printf("BitTorrent Download complete for %s\n", t.Name())
				app.EmitEvent("download_complete", "Download Finished!")
				return
			}
		}
	}()

	return t, nil
}
