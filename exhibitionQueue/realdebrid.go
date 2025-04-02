package exhibitionQueue

func (q *Queue) AddRealDebridDownloadToQueue(magnetLink string) {
	d := Download{
		Type:       RealDebridType,
		MagnetLink: magnetLink,
	}
	q.AddDownloadToQueue(d)
}
