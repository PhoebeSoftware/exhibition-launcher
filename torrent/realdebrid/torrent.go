package realdebrid

import (
	"fmt"
	"net/http"
	"net/url"
)

type AddMagnetResponse struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
}

type Torrent struct {
	ID       string   `json:"id"`
	Filename string   `json:"filename"`
	Hash     string   `json:"hash"`
	Bytes    int      `json:"bytes"`
	Host     string   `json:"host"`
	Split    int      `json:"split"`
	Progress int      `json:"progress"`
	Status   string   `json:"status"`
	Added    string   `json:"added"`
	Links    []string `json:"links"`
	Ended    *string  `json:"ended,omitempty"`
	Speed    *int     `json:"speed,omitempty"`
	Seeders  *int     `json:"seeders,omitempty"`
}

func (client *RealDebridClient) AddTorrentByMagnet(magnetLink string) (AddMagnetResponse, error) {
	params := url.Values{}
	params.Add("magnet", magnetLink)
	var result AddMagnetResponse
	req, err := client.newRequest(http.MethodPost, "/torrents/addMagnet", nil, params)
	if err != nil {
		return result, fmt.Errorf("could not add torrent by magnet: %w", err)
	}

	err = client.do(req, &result)
	if err != nil {
		return result, fmt.Errorf("could unmarshal response: %w", err)
	}

	return result, nil
}

func (client *RealDebridClient) GetTorrentInfoById(id string) (Torrent, error) {
	var result Torrent
	path := "/torrents/info/" + id
	fmt.Println(path)
	req, err := client.newRequest(http.MethodGet, path, nil, nil)

	if err != nil {
		return result, fmt.Errorf("could not get info from torrent: %w", err)
	}

	err = client.do(req, &result)
	if err != nil {
		return result, fmt.Errorf("could unmarshal response: %w", err)
	}

	return result, nil
}

func (client *RealDebridClient) GetTorents() ([]Torrent, error) {

	req, err := client.newRequest(http.MethodGet, "/torrents", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get request failed while requesting downloads: %w", err)
	}
	var result []Torrent

	err = client.do(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil

}
