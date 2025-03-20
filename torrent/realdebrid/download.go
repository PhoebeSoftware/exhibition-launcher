package realdebrid

import (
	"derpy-launcher072/utils/settings"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type DownloadItem struct {
	Id        string
	FileName  string
	MimeType  string
	FileSize  int
	Link      string
	Host      string
	Chunks    int
	Download  string
	Generated string
}

func (client *RealDebridClient) GetDownloads() ([]DownloadItem, error) {

	req, err := client.newRequest(http.MethodGet, "/downloads", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get request failed while requesting downloads: %w", err)
	}
	var result []DownloadItem

	err = client.do(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (client *RealDebridClient) DownloadByRDLink(link string, filePath string) error {
	req, err := http.NewRequest(http.MethodGet, link, nil)

	startTime := time.Now()

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get stats from file: %w", err)
	}

	rangeHeader := fmt.Sprintf("bytes=%d-", stat.Size())
	req.Header.Set("Range", rangeHeader)

	sizeResp, err := client.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not reach real debrid: %w", err)
	}
	fmt.Println(sizeResp.ContentLength)

	for i := int64(0); i < sizeResp.ContentLength; i += 10000000 {
		rangeHeader := fmt.Sprintf("bytes=%d-", stat.Size()+i)
		req.Header.Set("Range", rangeHeader)
		fmt.Println(i)

		resp, err := client.client.Do(req)
		if err != nil {
			return fmt.Errorf("could not encode link: %w", err)
		}

		defer resp.Body.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("could not copy files: %w", err)
		}

		fmt.Println("DONE")
	}

	fmt.Println("Done with: " + link)
	fmt.Println("Took: " + time.Since(startTime).String())
	return nil
}

func (client *RealDebridClient) DownloadByMagnet(magnetLink string, settings settings.Settings) error {
	addMagnetResponse, err := client.AddTorrentByMagnet(magnetLink)
	if err != nil {
		return err
	}

	torrent, err := client.GetTorrentInfoById(addMagnetResponse.Id)
	if err != nil {
		return err
	}

	err = client.SelectFiles(torrent)
	if err != nil {
		return err
	}

	// Re fetch torrent because torrent should now have selected files and links
	torrent, err = client.GetTorrentInfoById(addMagnetResponse.Id)
	if err != nil {
		return err
	}

	for _, link := range torrent.Links {
		unrestrictLink, err := client.UnrestrictLink(link)
		downloadPath := filepath.Join(settings.DownloadPath, unrestrictLink.Filename)
		if err != nil {
			return err
		}
		fmt.Println(unrestrictLink.Link)

		err = client.DownloadByRDLink(unrestrictLink.Download, downloadPath)
		if err != nil {
			return err
		}
	}

	return nil
}
