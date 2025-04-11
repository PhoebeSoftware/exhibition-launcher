package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	Cache map[string]Provider
}

var ProviderCacheDir = filepath.Join("cache", "providers")

func NewProviderManager() *ProviderManager {
	manager := &ProviderManager{
		Cache: map[string]Provider{},
	}

	err := os.MkdirAll(ProviderCacheDir, os.ModePerm)
	if err != nil {
		fmt.Println("could not create provider cache directory")
		return nil
	}

	LoadCachedToMemory(manager)
	return manager
}

func LoadCachedToMemory(p *ProviderManager) {
	entries, err := os.ReadDir(ProviderCacheDir)
	if err != nil {
		fmt.Printf("could not read provider cache directory: %v\n", err)
		return
	}

	for _, providerFile := range entries {
		if providerFile.IsDir() {
			continue
		}

		originalName := providerFile.Name()
		providerName := originalName[:len(originalName)-len(filepath.Ext(originalName))]

		_, ok := p.Cache[providerName]
		if ok {
			continue
		}

		filePath := filepath.Join(ProviderCacheDir, providerFile.Name())
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("could not open provider file %s: %v\n", filePath, err)
			continue
		}
		defer file.Close()

		var provider Provider
		err = json.NewDecoder(file).Decode(&provider)
		if err != nil {
			fmt.Printf("could not decode provider data from file %s: %v\n", filePath, err)
			continue
		}

		p.Cache[provider.ProviderName] = provider
	}
}

func IsSourceCached(link string) bool {
	file, err := os.Open(filepath.Join(ProviderCacheDir, filepath.Base(link)))
	if err != nil {
		fmt.Println("could not find source file")
		return false
	}
	defer file.Close()

	return true
}

func (p *ProviderManager) SearchDownloadsByGameName(query string) []ProviderDownload {
	var providerDownloads []ProviderDownload

	fmt.Println("starting search for", query)
	for _, provider := range p.Cache {
		for _, download := range provider.Downloads {
			go func() {
				if strings.Contains(strings.ToLower(download.Name), strings.ToLower(query)) {
					fmt.Println("found download", download.Name)
					providerDownloads = append(providerDownloads, download)
				}
			}()
		}
	}

	return providerDownloads
}

func (p *ProviderManager) CacheProvider(link string) error {
	lil := filepath.Base(link)
	providerName := lil[:len(lil)-len(filepath.Ext(lil))]

	fmt.Println(providerName)

	_, ok := p.Cache[providerName]
	if IsSourceCached(link) || ok {
		return errors.New("provider already cached")
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

	file, err := os.Create(filepath.Join(ProviderCacheDir, filepath.Base(link)))
	if err != nil {
		fmt.Println("could not cache provider data")
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
	p.Cache[data.ProviderName] = data

	fmt.Println("loaded source data from API")
	return nil
}
