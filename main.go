package main

import (
	"net/http"
	"time"

	"api.gotwitch.tk/routers"
	"api.gotwitch.tk/settings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	println("Loading settings...")
	settings.Setup()

	println("Loading database...")
	var err error
	settings.DB, err = gorm.Open(postgres.Open(settings.ServerSettings.PostgresDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	println("Finished initializing.")
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
