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

func SaveSource(source Source, link string) error {
	fmt.Println(filepath.Base(link))
	os.Mkdir("sources", os.ModePerm)

	file, err := os.Create(filepath.Join("sources", filepath.Base(link)))
	if err != nil {
		fmt.Println("could not save source data")
		return err
	}
	defer file.Close()

	bytes, err := json.Marshal(source)
	if err != nil {
		fmt.Println("could not marshal source data")
		return err
	}

	file.Write(bytes)
	return nil
}

func LoadSource(link string) (*Source, error) {
	fmt.Println(filepath.Base(link))

	file, err := os.Open(filepath.Join("sources", filepath.Base(link)))
	if err != nil {
		fmt.Println("could not find source file")
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

func (manager Manager) GetSource(link string) (*Source, error) {
	source, err := LoadSource(link)
	if err == nil {
		return source, nil
	}

	res, err := http.Get(link)
	if err != nil {
		fmt.Println("could get source data")
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)

	var data Source

	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("could decode source data")
		return nil, err
	}

	SaveSource(data, link)

	fmt.Println("loaded source data from API")
	return &data, nil
}
