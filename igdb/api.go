package igdb

import (
	"bytes"
	"encoding/json"
	"exhibition-launcher/utils/jsonUtils"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Image struct {
	ImageID string `json:"image_id"`
}

type ApiGame struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"summary"`
	Cover             Image  `json:"cover"`
	CoverURL          string
	Artworks          []Image `json:"artworks"`
	ArtworkUrlList    []string
	Screenshots       []Image `json:"screenshots"`
	ScreenshotUrlList []string
}

type APIManager struct {
	client   *http.Client
	settings *jsonModels.Settings
}



func (a *APIManager) SetupHeader(request *http.Request) {
	request.Header.Set("Client-ID", a.settings.IgdbSettings.IgdbClient)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.settings.IgdbSettings.IgdbAuth))
}

func NewAPI(settings *jsonModels.Settings, settingsManager *jsonUtils.JsonManager) (*APIManager, error) {
	apiManager := &APIManager{
		client:   &http.Client{},
		settings: settings,
	}

	// Generate new auth token if needed
	if time.Now().After(settings.IgdbSettings.ExpiresAt) {
		fmt.Println("Generating new auth token because the old one has expired")
		_, err := apiManager.GetAndSetNewAuthToken()
		if err != nil {
			fmt.Println("error fetching acces token:", err)
			return nil, err
		}
		// Save new auth token
		err = settingsManager.Save()
		if err != nil {
			return nil, err
		}
	}

	return apiManager, nil
}

func (a *APIManager) GetGameData(id int) (ApiGame, error) {
	header := fmt.Sprintf(`fields id, name, summary, cover.*, artworks.*, screenshots.*; where id = %d;`, id)

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return ApiGame{}, err
	}

	a.SetupHeader(request)

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

	a.SetupHeader(request)

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
