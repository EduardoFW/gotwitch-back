package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"api.gotwitch.tk/models"
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
	After      string
	Before     string
	First      int // max 100
	Game_id    string
	Language   string
	User_id    string
	User_login string
}

func GetStreamList(token string, clientId string, params *GetStreamListParams) (*models.StreamResponse, error) {
	var streamResponse models.StreamResponse
	client := http.Client{}

	// Default params
	if params.First == 0 {
		params.First = 100
	}

	// Transform struct to map
	paramsMap := make(map[string]string)
	paramsMap["after"] = params.After
	paramsMap["before"] = params.Before
	paramsMap["first"] = strconv.Itoa(params.First)
	paramsMap["game_id"] = params.Game_id
	paramsMap["language"] = params.Language
	paramsMap["user_id"] = params.User_id
	paramsMap["user_login"] = params.User_login

	// Build query string
	queryString := ""
	for key, value := range paramsMap {
		if value != "" {
			if queryString == "" {
				queryString += "?"
			} else {
				queryString += "&"
			}
			queryString += key + "=" + value
		}
	}

	url := "https://api.twitch.tv/helix/streams" + queryString

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
