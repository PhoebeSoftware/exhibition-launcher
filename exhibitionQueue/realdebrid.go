package exhibitionQueue

import (
	"exhibition-launcher/torrent/realdebrid"
	"fmt"
	"time"
)

func (q *Queue) AddRealDebridDownloadToQueue(magnetLink string) {
	d := Download{
		Type:             RealDebridType,
		MagnetLink:       magnetLink,
		DownloadProgress: realdebrid.DownloadProgress{},
	}
	q.AddDownloadToQueue(d)
}

func (q *Queue) StartRealDebridDownload(d Download) error {
	go func() {
		err := q.RealDebridClient.DownloadByMagnet(d.MagnetLink, q.DownloadPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	// Keep updating download progress
	for {
		d.DownloadProgress = q.RealDebridClient.DownloadProgress
		if d.DownloadProgress.IsDownloading {
			q.DownloadsInQueue[0] = d
			fmt.Println(d.DownloadProgress)
			<-time.NewTimer(10 * time.Second).C
		} else {
			break
		}
	}

	return nil
}


