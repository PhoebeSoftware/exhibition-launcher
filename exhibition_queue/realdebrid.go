package exhibition_queue

import "github.com/google/uuid"

func (q *Queue) AddRealDebridDownloadToQueue(magnetLink string) {
	d := Download{
		UUID:       uuid.New().String(),
		Type:       RealDebridType,
		MagnetLink: magnetLink,
	}
	q.AddDownloadToQueue(d)
}
