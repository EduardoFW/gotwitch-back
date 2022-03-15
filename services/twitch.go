package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"api.gotwitch.tk/models"
	"github.com/google/go-querystring/query"
)

func GetTwitchToken(clientID string, clientSecret string) (*models.TwitchToken, error) {
	var token models.TwitchToken

	resp, err := http.Post("https://id.twitch.tv/oauth2/token?client_id="+clientID+"&client_secret="+clientSecret+"&grant_type=client_credentials", "application/json", nil)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if token.AccessToken == "" {
		return nil, errors.New("No access token returned")
	}

	return &token, nil
}

type GetStreamListParams struct {
	After      string   `url:"after"`
	Before     string   `url:"before"`
	First      int      `url:"first"`
	Game_id    []string `url:"game_id"`
	Language   []string `url:"language"`
	User_id    []string `url:"user_id"`
	User_login []string `url:"user_login"`
}

func GetStreamList(token string, clientId string, params *GetStreamListParams) (*models.StreamResponse, error) {
	var streamResponse models.StreamResponse
	client := http.Client{}

	if params.First == 0 {
		params.First = 100
	}

	values, _ := query.Values(params)

	url := "https://api.twitch.tv/helix/streams?" + values.Encode()

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
