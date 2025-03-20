package realdebrid

import (
	"derpy-launcher072/utils/settingsManager"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func (client *RealDebridClient) DownloadByLink(link string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	defer file.Close()

	resp, err := client.client.Get(link)
	if err != nil {
		return fmt.Errorf("could not encode link: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("could not copy files: %w", err)
	}

	return nil
}

func (client *RealDebridClient) DownloadByMagnet(magnetLink string, settings *settingsManager.Settings) error {
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

	downloadPath := filepath.Join(settings.DownloadPath, "test.zip")

	for _, link := range torrent.Links {
		unrestrictLink, err := client.UnrestrictLink(link)
		if err != nil {
			return err
		}

		err = client.DownloadByLink(unrestrictLink.Link, downloadPath)
		if err != nil {
			return err
		}
		
	}

	
	return nil
}
