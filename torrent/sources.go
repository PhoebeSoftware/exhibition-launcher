package torrent

import (
	"encoding/json"
	"log"
	"net/http"
)

type SourceDownload struct {
	Name    string   `json:"title"`
	Magnets []string `json:"uris"`
	Date    string   `json:"uploadDate"`
	Size    string   `json:"fileSize"`
}

type Source struct {
	SourceName string           `json:"name"`
	Downloads  []SourceDownload `json:"downloads"`
}

func (manager Manager) GetSource(link string) Source {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Println(err)
	}

	res, err := manager.httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	decoder := json.NewDecoder(res.Body)

	var data Source
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err)
	}
	return data
}
