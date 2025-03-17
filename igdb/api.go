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
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Cover int    `json:"cover"`
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

func (*APIManager) GetCover(cover int) []Image {
	client := &http.Client{}

	header := fmt.Sprintf(`fields url; where id = %d;`, cover)

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/covers/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		fmt.Println(err)
		return []Image{}
	}

	SetupHeader(request)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return []Image{}
	}
	defer response.Body.Close()

	var images []Image
	jsonErr := json.NewDecoder(response.Body).Decode(&images)
	if jsonErr != nil {
		fmt.Println(err)
		return []Image{}
	}

	return images
}

func (*APIManager) GetGames(header string) []ApiGame {
	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.igdb.com/v4/games/", bytes.NewBuffer([]byte(header)))
	if err != nil {
		return nil
	}

	SetupHeader(request)

	response, err := client.Do(request)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	var games []ApiGame
	jsonErr := json.NewDecoder(response.Body).Decode(&games)
	if jsonErr != nil {
		return nil
	}

	fmt.Println(games)
	return games
}
