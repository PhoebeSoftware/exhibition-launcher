package realdebrid

import (
	"fmt"
	"net/http"
)

func (client *RealDebridClient) AddTorrent() {

}

type AvailableHost struct {
	Host        string `json:"host"`
	MaxFileSize int    `json:"max_file_size"`
}

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

type TraficInfo struct {
	Left  int
	Bytes int
	Links int
	Limit int
	Type  string
	Extra int
	Reset string
}

func (client *RealDebridClient) AvailableHosts() ([]AvailableHost, error) {
	req, err := client.newRequest(http.MethodGet, "/torrents/availableHosts", nil, "", nil)
	if err != nil {
		return nil, fmt.Errorf("get request failed while requesting available hosts: %w", err)
	}
	var result []AvailableHost

	err = client.do(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (client *RealDebridClient) GetDownloads() ([]DownloadItem, error) {

	req, err := client.newRequest(http.MethodGet, "/downloads", nil, "", nil)
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

func (client *RealDebridClient) GetTrafic() (map[string]TraficInfo, error) {

	type TraficResponse map[string]TraficInfo

	var traficResponse TraficResponse

	req, err := client.newRequest(http.MethodGet, "/traffic", nil, "", nil)
	if err != nil {
		return nil, fmt.Errorf("get request failed while requesting trafic: %w", err)
	}

	err = client.do(req, &traficResponse)
	if err != nil {
		return nil, err
	}

	return traficResponse, nil
}
