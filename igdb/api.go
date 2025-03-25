package igdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

type ApiGame struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"summary"`
	CoverID     int    `json:"cover"`
	MainCover string
}

type APIManager struct {
	client *http.Client
}

func SetupHeader(request *http.Request) {
	request.Header.Set("Client-ID", os.Getenv("IGDB_CLIENT"))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("IGDB_AUTH")))
}

// Special error so we can check in main
var (
	ErrorNoCoversFound = errors.New("could not find a cover with this id")
)

func NewAPI() *APIManager {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return &APIManager{client: &http.Client{}}
}

func (a *APIManager) GetCovers(coverID int) (string, error) {
	header := fmt.Sprintf(`fields image_id; where id = %d;`, coverID)
	var result string

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/covers/", bytes.NewBuffer([]byte(header)))
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
		fmt.Printf("No covers found with ID %d\n", coverID)
		return "", nil
	}
	imageID := images[0].ImageID

	imageURL := fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", imageID)
	return imageURL, nil
}

func (a *APIManager) GetGameData(id int) (ApiGame, error) {
	header := fmt.Sprintf(`fields id, name, summary, cover; where id = %d;`, id)

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return ApiGame{}, err
	}

	SetupHeader(request)

	response, err := a.client.Do(request)
	if err != nil {
		return ApiGame{}, err
	}
	defer response.Body.Close()

	var gameDataList []ApiGame
	jsonErr := json.NewDecoder(response.Body).Decode(&gameDataList)
	if jsonErr != nil {
		return ApiGame{}, err
	}

	if len(gameDataList) == 0 {
		return ApiGame{}, fmt.Errorf("no games found with id %d", id)
	}

	firstGameData := gameDataList[0]
	imageUrl, err := a.GetCovers(firstGameData.CoverID)
	if err != nil {
		// return game without cover
		return firstGameData, err
	}
	firstGameData.MainCover = imageUrl

	return firstGameData, nil
}

func (a *APIManager) GetGames(query string) []ApiGame {
	header := fmt.Sprintf(`fields id, name, summary, cover; search "%s";`, query)

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return []ApiGame{}
	}

	SetupHeader(request)

	response, err := a.client.Do(request)
	if err != nil {
		return []ApiGame{}
	}
	defer response.Body.Close()

	var games []ApiGame
	jsonErr := json.NewDecoder(response.Body).Decode(&games)
	if jsonErr != nil {
		return []ApiGame{}
	}
	return games
}
