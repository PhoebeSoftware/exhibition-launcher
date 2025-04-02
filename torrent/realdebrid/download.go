package realdebrid

import (
	"exhibition-launcher/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
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
func (client *RealDebridClient) DownloadDirectLink(link string, filePath string) error {
	startTime := time.Now()

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get stats from file: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	// A early request to fetch for the size of the file we are looping through later
	resp, err := client.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not reach real debrid: %w", err)
	}

	defer resp.Body.Close()
	if stat.Size() == resp.ContentLength {
		log.Printf("%v is already installed\n", filePath)
		return nil
	}

	totalSize := resp.ContentLength
	// 10mb
	sizeOfChunk := int64(10000000)

	fmt.Printf("Total file size: %d bytes\n", totalSize)

	var wg sync.WaitGroup
	var fileMutex sync.Mutex
	var downloadedBytes = stat.Size()

	numWorkers := client.Settings.RealDebridSettings.NumberOfThreads
	stopCh := make(chan interface{})
	errCh := make(chan error, 10)
	chunks := make(chan [2]int64, numWorkers)

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			case chunk, ok := <-chunks:
				if !ok {
					return
				}
				rangeStart, rangeEnd := chunk[0], chunk[1]
				req, err := http.NewRequest(http.MethodGet, link, nil)
				if err != nil {
					errCh <- fmt.Errorf("could not create request: %w", err)
					close(stopCh)
					return
				}

				rangeHeader := fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd)
				req.Header.Set("Range", rangeHeader)
				resp, err := client.client.Do(req)
				if err != nil {
					errCh <- fmt.Errorf("could not encode link: %w", err)
					close(stopCh)
					return
				}

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					errCh <- fmt.Errorf("could not read request: %w", err)
					close(stopCh)
					return
				}

				fileMutex.Lock()
				_, err = file.WriteAt(body, rangeStart)
				fileMutex.Unlock()
				if err != nil {
					errCh <- fmt.Errorf("could not copy files: %w", err)
					close(stopCh)
					return
				}

				atomic.AddInt64(&downloadedBytes, int64(len(body)))
			}
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		client.DownloadProgress.TotalBytes = resp.ContentLength
		client.DownloadProgress.IsDownloading = true
		for {
			select {
			case <-stopCh:
				client.DownloadProgress = DownloadProgress{IsDownloading: false}
				return
			default:
				//percent := float64(atomic.LoadInt64(&downloadedBytes)) / float64(totalSize) * 100
				client.DownloadProgress.DownloadedBytes = downloadedBytes
				client.DownloadProgress.Percent = (float64(downloadedBytes) / float64(resp.ContentLength)) * 100
			}
		}
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	for i := stat.Size(); i < totalSize; i += sizeOfChunk {
		rangeStart := i
		rangeEnd := rangeStart + sizeOfChunk + 1
		if rangeEnd >= totalSize {
			rangeEnd = totalSize - 1
		}
		chunks <- [2]int64{rangeStart, rangeEnd}
	}

	close(chunks)

	wg.Wait()

	close(errCh)

	close(stopCh)
	<-done

	for err := range errCh {
		return err
	}

	fmt.Println("\nDone with: " + filePath)
	fmt.Println("Took: " + time.Since(startTime).String())
	return nil
}

func (client *RealDebridClient) DownloadByMagnet(magnetLink string, path string) error {
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

	var unrestrictResponseList []UnrestrictResponse

	for _, link := range torrent.Links {
		unrestrictResponse, err := client.UnrestrictLink(link)
		if err != nil {
			return err
		}
		check, err := client.UnrestrictCheck(link)
		if err != nil {
			return err
		}

		if check != true {
			return Error503
		}

		unrestrictResponseList = append(unrestrictResponseList, unrestrictResponse)
	}

	if len(unrestrictResponseList) <= 0 {
		return ErrorNoLinksFound
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	disk := utils.DiskUsage(path)
	totalSize, err := client.GetDiskSizeOfAllLinks(unrestrictResponseList)
	if err != nil {
		return err
	}

	if int64(disk.Free) < totalSize {
		return ErrorNotEnoughDiskSpace
	}

	for _, unrestrictResponse := range unrestrictResponseList {
		downloadPath := filepath.Join(path, unrestrictResponse.Filename)
		fmt.Println(unrestrictResponse.Link)
		err = client.DownloadDirectLink(unrestrictResponse.Download, downloadPath)
		if err != nil {
			return err
		}
	}

	return nil
}
