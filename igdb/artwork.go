package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *APIManager) GetArtworkURL(artworkID int) (string, error) {
	header := fmt.Sprintf(`fields image_id; where id = %d;`, artworkID)
	var result string

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/artworks/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return result, err
	}

	SetupHeader(request)

	response, err := a.client.Do(request)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	var images []struct {
		ImageID string `json:"image_id"`
	}

	jsonErr := json.NewDecoder(response.Body).Decode(&images)
	if jsonErr != nil {
		return result, err
	}

	if len(images) == 0 {
		fmt.Printf("No covers found with ID %d\n", artworkID)
		return "", nil
	}
	imageID := images[0].ImageID

	imageURL := fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_1080p/%s.jpg", imageID)
	return imageURL, nil
}