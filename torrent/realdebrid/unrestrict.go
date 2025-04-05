package realdebrid

import (
	"fmt"
	"net/http"
	"net/url"
)

type UnrestrictResponse struct {
	ID         string `json:"id"`
	Filename   string `json:"filename"`
	MimeType   string `json:"mimeType"`
	Filesize   int64    `json:"filesize"`
	Link       string `json:"link"`
	Host       string `json:"host"`
	Chunks     int    `json:"chunks"`
	CRC        int    `json:"crc"`
	Download   string `json:"download"`
	Streamable int    `json:"streamable"`
}

func (client *RealDebridClient) UnrestrictLink(link string) (UnrestrictResponse, error) {
	var result UnrestrictResponse
	params := url.Values{}
	params.Add("link", link)

	req, err := client.newRequest(http.MethodPost, "/unrestrict/link", nil, params)
	if err != nil {
		return result, fmt.Errorf("could not encode link: %w", err)
	}

	err = client.do(req, &result)
	if err != nil {
		return result, fmt.Errorf("could not post link: %w", err)
	}

	return result, nil
}

type CheckResponse struct {
	Host      string `json:"host"`
	Link      string `json:"link"`
	Filename  string `json:"filename"`
	Filesize  int    `json:"filesize"`
	Supported int    `json:"supported"`
}

func (client *RealDebridClient) UnrestrictCheck(link string) (bool, error) {
	params := url.Values{}
	params.Add("link", link)

	req, err := client.newRequest(http.MethodPost, "/unrestrict/check", nil, params)
	if err != nil {
		return false, fmt.Errorf("could not encode link: %w", err)
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("could not post link: %w", err)
	}

	if resp.StatusCode == 503 {
		return false, nil
	}

	return true, nil
}
