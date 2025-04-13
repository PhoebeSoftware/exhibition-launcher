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
	ProviderURL  string             `json:"url"`
	ETag         string             `json:"etag"`
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

	VerifyAllLocalProviders(manager)
	fmt.Println("verified all providers")
	LoadLocalToMemory(manager)
	fmt.Println("loaded all providers")
	return manager
}

func VerifyAllLocalProviders(p *ProviderManager) {
	entries, err := os.ReadDir(ProviderDir)
	if err != nil {
		fmt.Printf("could not read provider directory: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	mutex := sync.Mutex{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		wg.Add(1)

		go func(providerFile os.DirEntry) {
			defer wg.Done()

			file, err := os.Open(filepath.Join(ProviderDir, providerFile.Name()))
			if err != nil {
				fmt.Printf("couldnt create provider data for %s: %v\n", providerFile.Name(), err)
				return
			}
			defer file.Close()

			var serialisedProvider Provider

			err = json.NewDecoder(file).Decode(&serialisedProvider)
			if err != nil {
				fmt.Println("could not marshal provider data")
				return
			}

			mutex.Lock()
			err = p.VerifyProvider(serialisedProvider, providerFile.Name())
			mutex.Unlock()

			if err != nil {
				fmt.Printf("failed to verify provider %s: %v\n", providerFile.Name(), err)
				return
			}
		}(entry)
	}

	wg.Wait()
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

			// provider already laoded
			_, exists := p.Providers[providerName]
			if exists {
				return
			}

			// open provider file
			filePath := filepath.Join(ProviderDir, providerFile.Name())
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("could not open provider file %s: %v\n", filePath, err)
				return
			}
			defer file.Close()

			// read provider data
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

// invalid or valid
func (p *ProviderManager) VerifyProvider(provider Provider, providerFile string) error {
	fmt.Println("verifying provider", providerFile)

	res, err := http.Head(provider.ProviderURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fmt.Printf("remote status code: %d\n", res.StatusCode)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("provider %s is not available", provider.ProviderName)
	}

	fmt.Printf("Remote ETag: %s, Local ETag: %s\n", res.Header.Get("ETag"), provider.ETag)
	if provider.ETag == res.Header.Get("ETag") {
		fmt.Printf("Provider %s is up to date\n", provider.ProviderName)
		return nil
	}

	fmt.Printf("Provider %s is outdated\n", provider.ProviderName)

	// remove local provider
	err = os.Remove(filepath.Join(ProviderDir, providerFile))
	if err != nil {
		return err
	}

	// download new provider
	err = p.DownloadProvider(provider.ProviderURL)
	if err != nil {
		return err
	}

	return nil
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
	var mutex sync.Mutex
	var wg sync.WaitGroup

	fmt.Println("starting search for", query)
	for _, provider := range p.Providers {
		for _, download := range provider.Downloads {
			wg.Add(1)
			go func(providerName string, download ProviderDownload) {
				defer wg.Done()
				if strings.Contains(strings.ToLower(download.Name), strings.ToLower(query)) {
					fmt.Println("found download", download.Name)
					mutex.Lock()
					providerDownloads[providerName] = download
					mutex.Unlock()
				}
			}(provider.ProviderName, download)
		}
	}

	wg.Wait()
	return providerDownloads
}

func (p *ProviderManager) DownloadProvider(link string) error {
	providerBase := filepath.Base(link)
	providerName := providerBase[:len(providerBase)-len(filepath.Ext(providerBase))]

	_, ok := p.Providers[providerName]
	if IsProviderDownloaded(link) || ok {
		return nil
	}

	res, err := http.Get(link)
	if err != nil {
		return errors.New("couldnt get provider data")
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("provider with link %s is not available", link)
	}

	decoder := json.NewDecoder(res.Body)

	var data Provider

	err = decoder.Decode(&data)
	if err != nil {
		return errors.New("couldnt decode provider data")
	}
	data.ETag = res.Header.Get("ETag")
	data.ProviderURL = link

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
