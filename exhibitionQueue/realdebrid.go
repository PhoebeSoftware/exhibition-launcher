package exhibitionQueue

import (
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func (q *Queue) AddRealDebridDownloadToQueue(magnetLink string) {
	d := Download{
		Type:       RealDebridType,
		MagnetLink: magnetLink,
	}
	q.AddDownloadToQueue(d)
}

func (q *Queue) StartRealDebridDownload(app *application.App, d Download) error {
	go func() {
		err := q.RealDebridClient.DownloadByMagnet(app, d.MagnetLink, q.DownloadPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	return nil
}
