package queue

import (
	"exhibition-launcher/torrent"
	"exhibition-launcher/torrent/realdebrid"
)

var (
	RealDebridType = "Real-Debrid"
	TorrentType = "Torrent"
)

type Queue struct {
	DownloadsInQueue map[int]Download
	TorrentManager torrent.Manager
	RealDebridClient realdebrid.RealDebridClient
	Path string
}


type Download struct {
	Type string
	MagnetLink string
	DownloadProgress realdebrid.DownloadProgress
}

func (q *Queue) StartRealDebridDownload(d Download) {
	q.RealDebridClient.DownloadByMagnet(d.MagnetLink, )
}

func (q *Queue) AddDownloadToQueue(d Download) {
	priority := len(q.DownloadsInQueue) + 1
	q.DownloadsInQueue[priority] = d
}

func (q *Queue) AddRealDebridDownloadToQueue(magnetLink string) {
	priority := len(q.DownloadsInQueue) + 1
	d := Download{
		Type:             RealDebridType,
		MagnetLink:       magnetLink,
		DownloadProgress: realdebrid.DownloadProgress{},
	}
	q.DownloadsInQueue[priority] = d
}