package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"api.gotwitch.tk/controllers/jobs"
	"api.gotwitch.tk/models"
	"api.gotwitch.tk/routers"
	"api.gotwitch.tk/settings"

	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         settings.ServerSettings.SentryDSN,
		Environment: settings.ServerSettings.SentryEnvironment,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func connectToDatabase() error {
	var err error
	settings.DB, err = gorm.Open(postgres.Open(settings.ServerSettings.PostgresDSN), &gorm.Config{})
	return err
}

func tryToConnectToDatabaseAndRetry(retryAttempts int) error {
	for i := 0; i < retryAttempts; i++ {
		if err := connectToDatabase(); err == nil {
			return nil
		}
		time.Sleep(time.Second * 5)
	}
	return errors.New("Could not connect to database")
}

func init() {
	println("Loading settings...")
	settings.DefaultSetup()

	println("Loading database...")
	if err := tryToConnectToDatabaseAndRetry(5); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	println("Adding database migrations...")
	settings.DB.AutoMigrate(&models.Stream{})
	settings.DB.AutoMigrate(&models.Job{})
	settings.DB.AutoMigrate(&models.Category{})

	println("Initializing scheduler...")
	go scheduler()

	initSentry()

	println("Finished initializing.")

}

func scheduler() {
	println("Starting scheduler...")
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(15).Minutes().SingletonMode().Do(jobs.Orchestrator)
	scheduler.Every(15).Days().SingletonMode().Do(jobs.CategoryOrchestrator)
	scheduler.StartBlocking()
}

func main() {
	defer sentry.Flush(2 * time.Second)

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
