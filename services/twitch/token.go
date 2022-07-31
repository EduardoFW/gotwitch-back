package twitch

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"api.gotwitch.tk/settings"
)

var lock = &sync.Mutex{}

type TwitchToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenManager struct {
	token       *TwitchToken
	expiesAfter time.Time
}

func (t *TokenManager) GetToken() (string, error) {

	if t.token == nil || time.Now().After(t.expiesAfter) {
		token, err := t.GetTwitchToken(settings.ServerSettings.TwitchClientID, settings.ServerSettings.TwitchClientSecret)

		if err != nil {
			return "", err
		}

		t.token = token
		t.expiesAfter = time.Now().Add(time.Second * time.Duration(t.token.ExpiresIn))
	}

	return t.token.AccessToken, nil
}

func (t *TokenManager) GetTwitchToken(clientID string, clientSecret string) (*TwitchToken, error) {
	var token TwitchToken

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

var singleInstanceTokenManager *TokenManager
var onceTokenManager sync.Once

func GetTokenManager() *TokenManager {
	onceTokenManager.Do(func() {
		singleInstanceTokenManager = &TokenManager{}
	})
	return singleInstanceTokenManager
}
