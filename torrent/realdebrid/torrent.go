package realdebrid

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AddMagnetResponse struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
}

type TorrentFile struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	Bytes    int    `json:"bytes"`
	Selected int    `json:"selected"`
}

type Torrent struct {
	ID               string        `json:"id"`
	Filename         string        `json:"filename"`
	OriginalFilename string        `json:"original_filename"`
	Hash             string        `json:"hash"`
	Bytes            int           `json:"bytes"`
	OriginalBytes    int           `json:"original_bytes"`
	Host             string        `json:"host"`
	Split            int           `json:"split"`
	Progress         int           `json:"progress"`
	Status           string        `json:"status"`
	Added            string        `json:"added"`
	TorrentFiles     []TorrentFile `json:"files"`
	Links            []string      `json:"links"`
	Ended            *string       `json:"ended,omitempty"`
	Speed            *int          `json:"speed,omitempty"`
	Seeders          *int          `json:"seeders,omitempty"`
}

func (client *RealDebridClient) AddTorrentByMagnet(magnetLink string) (AddMagnetResponse, error) {
	params := url.Values{}
	params.Add("magnet", magnetLink)
	var result AddMagnetResponse
	req, err := client.newRequest(http.MethodPost, "/torrents/addMagnet", nil, params)
	if err != nil {
		return result, fmt.Errorf("error while encoding url: %w", err)
	}

	err = client.do(req, &result)
	if err != nil {
		return result, fmt.Errorf("error while adding torrent by magnet link: %w", err)
	}

	fmt.Println(req.URL.String())
	return result, nil
}

func (client *RealDebridClient) GetTorrentInfoById(id string) (Torrent, error) {
	var result Torrent
	req, err := client.newRequest(http.MethodGet, "/torrents/info/"+id, nil, nil)

	if err != nil {
		return result, fmt.Errorf("error while encoding url: %w", err)
	}

	err = client.do(req, &result)
	if err != nil {
		return result, fmt.Errorf("error while requesting torrent by id: %w", err)
	}

	return result, nil
}

func (client *RealDebridClient) GetTorrents() ([]Torrent, error) {
	req, err := client.newRequest(http.MethodGet, "/torrents", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error while encoding url: %w", err)
	}
	var result []Torrent

	err = client.do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("error while requesting torrents: %w", err)
	}

	return result, nil

}

// SelectFiles This function is supposed to be called after AddTorrentByMagnet() or similar. Because Real debrid needs to know which files to torrent
func (client *RealDebridClient) SelectFiles(torrent Torrent) error {
	params := url.Values{}
	var fileIDs []string
	for _, file := range torrent.TorrentFiles {
		fileIDs = append(fileIDs, strconv.Itoa(file.ID))
	}
	filesParam := strings.Join(fileIDs, ",")
	params.Set("files", filesParam)

	req, err := client.newRequest(http.MethodPost, "/torrents/selectFiles/"+torrent.ID, nil, params)
	if err != nil {
		return fmt.Errorf("error while encoding url: %w", err)
	}
	err = client.do(req, nil)
	if err != nil {
		return fmt.Errorf("get request failed while posting select files: %w", err)
	}
	return nil
}

// CheckIfTorrentAlreadyExists Returns torrent id name is misleading
func (client *RealDebridClient) CheckIfTorrentAlreadyExists(magnetLink string) (string, error) {
	torrents, err := client.GetTorrents()
	if err != nil {
		return "", err
	}

	magnetHash, err := GetMagnetLinkHash(magnetLink)
	if err != nil {
		return "", err
	}

	for _, torrent := range torrents {
		if torrent.Hash == magnetHash {
			return torrent.ID, err
		}
	}

	return "", nil
}
