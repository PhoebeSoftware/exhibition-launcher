package exhibitionQueue

import "github.com/google/uuid"

func (q *Queue) AddTorrentDownloadToQueue(magnetLink string) {
	d := Download{
		UUID:       uuid.New().String(),
		Type:       TorrentType,
		MagnetLink: magnetLink,
	}
	q.AddDownloadToQueue(d)
}
