package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *APIManager) GetCover(coverID int) (string, error) {
	header := fmt.Sprintf(`fields image_id; where id = %d;`, coverID)
	var result string

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/covers/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return result, err
	}

	a.SetupHeader(request)

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
		fmt.Printf("No covers found with ID %d\n", coverID)
		return "", nil
	}
	// Grab the image ID of the first image out of the array
	imageID := images[0].ImageID

	result = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", imageID)
	return result, nil
}
