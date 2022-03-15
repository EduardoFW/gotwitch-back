package controllers

import (
	"sync"
	"time"

	"api.gotwitch.tk/models"
	"api.gotwitch.tk/services"
	"api.gotwitch.tk/settings"
)

var lock = &sync.Mutex{}

type TokenManager struct {
	token       *models.TwitchToken
	expiesAfter time.Time
}

func (t *TokenManager) GetToken() (string, error) {

	if t.token == nil || time.Now().After(t.expiesAfter) {
		token, err := services.GetTwitchToken(settings.ServerSettings.TwitchClientID, settings.ServerSettings.TwitchClientSecret)

		if err != nil {
			return "", err
		}

		t.token = token
		t.expiesAfter = time.Now().Add(time.Second * time.Duration(t.token.ExpiresIn))
	}

	return t.token.AccessToken, nil
}

var singleInstance *TokenManager
var once sync.Once

func GetTokenManager() *TokenManager {
	once.Do(func() {
		singleInstance = &TokenManager{}
	})
	return singleInstance
}
