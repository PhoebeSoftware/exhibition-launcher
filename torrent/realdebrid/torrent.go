package realdebrid

import (
	"fmt"
	"net/http"
	"net/url"
)

func (client *RealDebridClient) AddTorrentByMagnet(magnetLink string) error {
	params := url.Values{}
	params.Add("files", magnetLink)

	req, err := client.newRequest(http.MethodPost, "/torrents/addMagnet", nil, params, nil)
	if err != nil {
		return fmt.Errorf("could not add torrent by magnet: %w", err)
	}

	fmt.Println(req)
	return nil
}



