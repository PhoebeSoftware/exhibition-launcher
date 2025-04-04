package igdb

import (
	"bytes"
	"encoding/json"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (a *APIManager) GetAndSetNewAuthToken() (string, error) {
	client := a.client
	params := url.Values{}
	params.Add("client_id", a.settings.IgdbSettings.IgdbClient)
	params.Add("client_secret", a.settings.IgdbSettings.IgdbSecret)
	params.Add("grant_type", "client_credentials")
	uri := "https://id.twitch.tv/oauth2/token" + "?" + params.Encode()
	req, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		return "", fmt.Errorf("error setting up request:%w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	type AuthResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while requesting:%w", err)
	}

	var authResponse AuthResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body:%w", err)
	}

	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding json:%w", err)
	}

	a.settings.IgdbSettings.IgdbAuth = authResponse.AccessToken
	a.settings.IgdbSettings.ExpiresIn = authResponse.ExpiresIn

	return authResponse.AccessToken, nil
}

func (a *APIManager) SetupHeader(request *http.Request) {
	request.Header.Set("Client-ID", a.settings.IgdbSettings.IgdbClient)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.settings.IgdbSettings.IgdbAuth))
}

func NewAPI(settings *jsonModels.Settings) (*APIManager, error) {
	return &APIManager{
		client:   &http.Client{},
		settings: settings,
	}, nil
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
