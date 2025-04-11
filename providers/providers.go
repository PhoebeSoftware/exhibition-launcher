package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ProviderDownload struct {
	Name    string   `json:"title"`
	Magnets []string `json:"uris"`
	Date    string   `json:"uploadDate"`
	Size    string   `json:"fileSize"`
}

type Provider struct {
	ProviderName string             `json:"name"`
	Downloads    []ProviderDownload `json:"downloads"`
}

type ProviderManager struct {
	Providers map[string]Provider
}

var ProviderDir = filepath.Join("provider")

func NewProviderManager() *ProviderManager {
	manager := &ProviderManager{
		Providers: map[string]Provider{},
	}

	err := os.MkdirAll(ProviderDir, os.ModePerm)
	if err != nil {
		fmt.Println("could not create provider directory")
		return nil
	}

	LoadLocalToMemory(manager)
	return manager
}

func LoadLocalToMemory(p *ProviderManager) {
	entries, err := os.ReadDir(ProviderDir)
	if err != nil {
		fmt.Printf("could not read provider directory: %v\n", err)
		return
	}

	mutex := sync.Mutex{}

	var wg sync.WaitGroup

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		wg.Add(1)

		go func(providerFile os.DirEntry) {
			defer wg.Done()

			originalName := providerFile.Name()
			providerName := originalName[:len(originalName)-len(filepath.Ext(originalName))]

			_, exists := p.Providers[providerName]
			if exists {
				return
			}

			filePath := filepath.Join(ProviderDir, providerFile.Name())
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("could not open provider file %s: %v\n", filePath, err)
				return
			}
			defer file.Close()

			var provider Provider
			err = json.NewDecoder(file).Decode(&provider)
			if err != nil {
				fmt.Printf("could not decode provider data from file %s: %v\n", filePath, err)
				return
			}

			mutex.Lock()
			p.Providers[provider.ProviderName] = provider
			mutex.Unlock()
		}(entry)
	}

	wg.Wait()
}

func IsProviderDownloaded(link string) bool {
	file, err := os.Open(filepath.Join(ProviderDir, filepath.Base(link)))
	if err != nil {
		return false
	}
	defer file.Close()

	return true
}

func (p *ProviderManager) SearchDownloadsByGameName(query string) map[string]ProviderDownload {
	var providerDownloads = map[string]ProviderDownload{}

	fmt.Println("starting search for", query)
	for _, provider := range p.Providers {
		for _, download := range provider.Downloads {
			go func() {
				if strings.Contains(strings.ToLower(download.Name), strings.ToLower(query)) {
					fmt.Println("found download", download.Name)
					providerDownloads[provider.ProviderName] = download
				}
			}()
		}
	}

	return providerDownloads
}

func (p *ProviderManager) DownloadProvider(link string) error {
	lil := filepath.Base(link)
	providerName := lil[:len(lil)-len(filepath.Ext(lil))]

	_, ok := p.Providers[providerName]
	if IsProviderDownloaded(link) || ok {
		return nil
	}

	res, err := http.Get(link)
	if err != nil {
		return errors.New("couldnt get provider data")
	}

	decoder := json.NewDecoder(res.Body)

	var data Provider

	err = decoder.Decode(&data)
	if err != nil {
		return errors.New("couldnt decode provider data")
	}

	file, err := os.Create(filepath.Join(ProviderDir, filepath.Base(link)))
	if err != nil {
		fmt.Println("couldnt create provider data")
		return err
	}
	defer file.Close()

	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("could not marshal provider data")
		return err
	}

	// write to disk
	file.Write(bytes)

	// add to memory
	p.Providers[data.ProviderName] = data

	fmt.Println("loaded source data from API")
	return nil
}
