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

type QueueStatus int

const (
	Idle QueueStatus = iota
	Downloading
	Extracting
)

type Queue struct {
	DownloadsInQueue []Download
	TorrentManager   *torrent.Manager
	RealDebridClient *realdebrid.RealDebridClient
	App              *application.App

	DownloadPath string
	QueueStatus  QueueStatus
	Paused       bool
}

type Download struct {
	UUID       string
	Type       string
	MagnetLink string
	Progress   float64
}

func (q *Queue) GetFirstDownload() Download {
	return q.DownloadsInQueue[0]
}

func (q *Queue) SetStatus(status QueueStatus) {
	q.QueueStatus = status
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

func (q *Queue) SetPaused(value bool) {
	q.Paused = value

	switch q.GetFirstDownload().Type {
	case RealDebridType:
		q.RealDebridClient.SetPaused(value)
	case TorrentType:
		q.TorrentManager.SetPaused(value)
	}
}

func (q *Queue) GetPaused() bool {
	return q.Paused
}

// StartDownloads starts the first donwload from the q.DownloadsInQueue. When done removes the first one and queues again until there are no more items in the array.
func (q *Queue) StartDownloads() error {
	if q.QueueStatus != Idle {
		return fmt.Errorf("already downloading")
	}
	if len(q.DownloadsInQueue) <= 0 {
		return fmt.Errorf("no downloads in queue")
	}

	defer q.SetStatus(Idle)
	q.SetStatus(Downloading)

	download := q.GetFirstDownload()

	switch download.Type {
	case RealDebridType:
		fmt.Println("Starting Real-Debrid download!!")
		err := q.RealDebridClient.DownloadByMagnet(q.App, download.MagnetLink, q.DownloadPath)
		if err != nil {
			return err
		}
	case TorrentType:
		fmt.Println("Starting BitTorrent download!!")
		t, err := q.TorrentManager.AddTorrent(q.App, download.UUID, download.MagnetLink)
		if err != nil {
			return err
		}

		// wacht
		<-t.NotifyComplete()
	}

	fmt.Printf("%s download complete\n", download.Type)

	// Next download
	q.SetStatus(Idle)
	q.RemoveFromQueue(0)

	err := q.StartDownloads()
	if err != nil {
		return err
	}

	return nil
}
