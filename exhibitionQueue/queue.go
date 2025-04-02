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
	Downloading      bool
}

type Download struct {
	Type       string
	MagnetLink string
	Progress   float64
}

func (q *Queue) GetCurrentDownload() Download {
	return q.DownloadsInQueue[0]
}

func (q *Queue) GetDownloadInQueue() []Download {
	return q.DownloadsInQueue
}

func (q *Queue) AddDownloadToQueue(d Download) {
	q.DownloadsInQueue = append(q.DownloadsInQueue, d)
}

func (q *Queue) RemoveFromQueue(i int) {
	q.DownloadsInQueue = append(q.DownloadsInQueue[:i], q.DownloadsInQueue[i+1:]...)
}

func (q *Queue) SetDownloading(value bool) {
	q.Downloading = value
}
func (q *Queue) StartDownloads() error {
	if q.Downloading {
		return fmt.Errorf("already downloading")
	}
	if len(q.DownloadsInQueue) <= 0 {
		return fmt.Errorf("no downloads in queue")
	}
	defer q.SetDownloading(false)
	q.SetDownloading(true)
	download := q.DownloadsInQueue[0]
	switch download.Type {
	case RealDebridType:
		fmt.Println("Starting Real-Debrid download!!")
		err := q.RealDebridClient.DownloadByMagnet(q.App, download.MagnetLink, q.DownloadPath)
		if err != nil {
			return err
		}
	case TorrentType:
		// OMAR LOCK IN!!!
	}

	// Next download
	q.SetDownloading(false)
	q.DownloadsInQueue = append(q.DownloadsInQueue[:0], q.DownloadsInQueue[1:]...)
	err := q.StartDownloads()
	if err != nil {
		return err
	}
	return nil
}
