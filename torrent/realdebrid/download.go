package realdebrid

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
