package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *APIManager) GetArtworkURLs(artworkIDs []int) ([]string, error) {
	var result []string

	for _, id := range artworkIDs {
		fmt.Println(id)
		header := fmt.Sprintf(`fields image_id; where id = %d;`, id)
		request, err := http.NewRequest("POST", "https://api.igdb.com/v4/artworks/", bytes.NewBuffer([]byte(header)))
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
			fmt.Printf("No banners found with ID %d\n", id)
			return result, nil
		}

		for _, image := range images {
			imageID := image.ImageID
			imageURL := fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_1080p/%s.jpg", imageID)
			result = append(result, imageURL)
			fmt.Println("added " + imageURL)
		}
	}

	return result, nil
}
