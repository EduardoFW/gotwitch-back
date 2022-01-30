package main

import (
	"net/http"
	"time"

	"github.com/EduardoWeber/twitch-go/routers"
	"github.com/EduardoWeber/twitch-go/settings"
)

func init() {
	settings.Setup()
}

func main() {

	router := routers.InitRouter()
	s := &http.Server{
		Handler:        router,
		Addr:           settings.ServerSettings.Addr,
		ReadTimeout:    time.Duration(settings.ServerSettings.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(settings.ServerSettings.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(settings.ServerSettings.IdleTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
