package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"api.gotwitch.tk/models"
)

func GetTwitchToken(clientID string, clientSecret string) (string, error) {
	var token models.TwitchToken

	resp, err := http.Post("https://id.twitch.tv/oauth2/token?client_id="+clientID+"&client_secret="+clientSecret+"&grant_type=client_credentials", "application/json", nil)
	if err != nil {
		return "", err
	}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if token.AccessToken == "" {
		return "", errors.New("No access token returned")
	}

	return string(token.AccessToken), nil
}

func GetStreamList(token string, clientId, after string) (*models.StreamResponse, error) {
	var streamResponse models.StreamResponse
	client := http.Client{}

	// Querystring parameters
	params := "?first=100"
	if after != "" {
		params += "&after=" + after
	}

	url := "https://api.twitch.tv/helix/streams" + params

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
		"Client-Id":     []string{clientId},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&streamResponse)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return &streamResponse, nil
}
