package igdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Image struct {
	ImageID string `json:"image_id"`
}

type ApiGame struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"summary"`

	Cover             Image    `json:"-"`
	CoverURL          string   `json:"-"`
	Artworks          []Image  `json:"-"`
	ArtworkUrlList    []string `json:"-"`
	Screenshots       []Image  `json:"-"`
	ScreenshotUrlList []string `json:"-"`
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

func (a *APIManager) GetGameData(id int) (ApiGame, error) {
	header := fmt.Sprintf(`fields id, name, summary, cover.*, artworks.*, screenshots.*; where id = %d;`, id)

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
	fmt.Println(firstGameData.Name+" :", firstGameData.Id)
	imageID := firstGameData.Cover.ImageID
	imageURL := fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", imageID)
	firstGameData.CoverURL = imageURL
	fmt.Println("added cover " + imageURL)

	for _, image := range firstGameData.Artworks {
		imageID := image.ImageID
		imageURL := fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_1080p/%s.jpg", imageID)
		firstGameData.ArtworkUrlList = append(firstGameData.ArtworkUrlList, imageURL)
		fmt.Println("added artwork " + imageURL)
	}

	for _, image := range firstGameData.Screenshots {
		imageID := image.ImageID
		imageURL := fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_1080p/%s.jpg", imageID)
		firstGameData.ScreenshotUrlList = append(firstGameData.ScreenshotUrlList, imageURL)
		fmt.Println("added screenshot " + imageURL)
	}

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
