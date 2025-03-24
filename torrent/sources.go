package torrent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

func SaveSource(source Source, link string) {
	os.Mkdir("sources", os.ModePerm)

	file, err := os.Create(filepath.Join("sources", link+".json"))
	if err != nil {
		fmt.Println("could not save source data")
	}
	defer file.Close()

	bytes, err := json.Marshal(source)
	if err != nil {
		fmt.Println("could not marshal source data")
	}

	file.Write(bytes)
}

func LoadSource(link string) (*Source, error) {
	file, err := os.Open(filepath.Join("sources", link+".json"))
	if err != nil {
		fmt.Println("could not open source file")
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var data Source
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("could not decode source data")
		return nil, err
	}

	fmt.Println("loaded source data from file")
	return &data, nil
}

func (manager Manager) GetSource(link string) Source {
	linkFilename := filepath.Base(link)
	fmt.Println(linkFilename)

	source, err := LoadSource(linkFilename)
	if err == nil {
		return *source
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println("could get source data")
	}

	res, err := manager.httpClient.Do(req)
	if err != nil {
		fmt.Println("could request source data")
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var data Source
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("could decode source data")
	}

	SaveSource(data, link)

	fmt.Println("loaded source data from API")
	return data
}
