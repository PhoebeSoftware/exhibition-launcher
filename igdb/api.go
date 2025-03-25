package igdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ApiGame struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"summary"`
	Cover       int    `json:"cover"`
}

type Image struct {
	Link string `json:"url"`
}

type APIManager struct {
	client *http.Client
}

func SetupHeader(request *http.Request) {
	request.Header.Set("Client-ID", os.Getenv("IGDB_CLIENT"))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("IGDB_AUTH")))
}

func NewAPI() *APIManager {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return &APIManager{client: &http.Client{}}
}

func (a *APIManager) GetCover(cover int) []string {
	header := fmt.Sprintf(`fields url; where id = %d;`, cover)

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/covers/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	SetupHeader(request)

	response, err := a.client.Do(request)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	defer response.Body.Close()

	var images []Image
	jsonErr := json.NewDecoder(response.Body).Decode(&images)
	if jsonErr != nil {
		fmt.Println(err)
		return []string{}
	}

	urls := make([]string, len(images))
	for i, img := range images {
		urls[i] = img.Link
	}

	return urls
}

func (a *APIManager) GetGameData(id int) ApiGame {
	header := fmt.Sprintf(`fields id, name, summary, cover; where id = %d;`, id)

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return ApiGame{}
	}

	SetupHeader(request)

	response, err := a.client.Do(request)
	if err != nil {
		return ApiGame{}
	}
	defer response.Body.Close()

	var game []ApiGame
	jsonErr := json.NewDecoder(response.Body).Decode(&game)
	if jsonErr != nil {
		return ApiGame{}
	}

	if len(game) == 0 {
		return ApiGame{}
	}
	return game[0]
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
