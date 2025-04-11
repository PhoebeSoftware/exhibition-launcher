package exhibitionQueue

import "github.com/google/uuid"

func (q *Queue) AddTorrentDownloadToQueue(magnetLink string) {
	q.AddDownloadToQueue(Download{
		UUID:       uuid.New().String(),
		Type:       TorrentType,
		MagnetLink: magnetLink,
	})
}
