package realdebrid

import (
	"fmt"
	"net/http"
)

type AvailableHost struct {
	Host        string `json:"host"`
	MaxFileSize int    `json:"max_file_size"`
}

func (client *RealDebridClient) AvailableHosts() ([]AvailableHost, error) {
	req, err := client.newRequest(http.MethodGet, "/torrents/availableHosts", nil, nil)
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