package exhibitionQueue

import (
	"exhibition-launcher/torrent"
	"exhibition-launcher/torrent/realdebrid"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	RealDebridType = "Real-Debrid"
	TorrentType    = "Torrent"
)

type Queue struct {
	DownloadsInQueue []Download
	TorrentManager   *torrent.Manager
	RealDebridClient *realdebrid.RealDebridClient
	DownloadPath     string
	App              *application.App
}

type Download struct {
	Type       string
	MagnetLink string
	Progress   float64
}

func (q *Queue) GetCurrentDownload() Download {
	fmt.Println("Getting current download!!!!!!")
	return q.DownloadsInQueue[0]
}

func (q *Queue) AddDownloadToQueue(d Download) {
	q.DownloadsInQueue = append(q.DownloadsInQueue, d)
}

func (q *Queue) RemoveFromQueue(i int) {
	q.DownloadsInQueue = append(q.DownloadsInQueue[:i], q.DownloadsInQueue[i+1:]...)
}
func (q *Queue) StartDownloads() error {
	if len(q.DownloadsInQueue) <= 0 {
		fmt.Println("No downloads in queue")
		return nil
	}
	download := q.DownloadsInQueue[0]
	switch download.Type {
	case RealDebridType:
		fmt.Println("Starting Real-Debrid download!!")
		err := q.StartRealDebridDownload(q.App, download)
		if err != nil {
			return err
		}
	case TorrentType:
		// OMAR LOCK IN!!!
	}

	// Next download

	return nil
}
