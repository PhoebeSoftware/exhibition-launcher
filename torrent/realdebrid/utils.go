package realdebrid

import (
	"fmt"
	"net/http"
)

func (client *RealDebridClient) GetDiskSizeOfAllLinks(unrestrictResponses []UnrestrictResponse) (int64, error) {
	var (
		totalSize int64
	)

	for _, unrestrictResponse := range unrestrictResponses {
		req, err := http.NewRequest(http.MethodGet, unrestrictResponse.Download, nil)
		if err != nil {
			return totalSize, fmt.Errorf("could create request: %w", err)
		}

		resp, err := client.client.Do(req)
		if err != nil {
			return totalSize, fmt.Errorf("could not fetch size: %w", err)
		}

		totalSize += resp.ContentLength
	}

	return totalSize, nil
}
