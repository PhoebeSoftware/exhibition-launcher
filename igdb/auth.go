package igdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

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
	expiresAt := time.Now().Add(time.Duration(authResponse.ExpiresIn) * time.Second)
	a.settings.IgdbSettings.ExpiresAt = expiresAt
	return authResponse.AccessToken, nil
}
