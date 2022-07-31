package jobs

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"api.gotwitch.tk/models"
	"api.gotwitch.tk/services/twitch"
	"api.gotwitch.tk/settings"
	"gorm.io/gorm"
)

func UpsertStreams(streams []models.Stream, job models.Job) error {
	return settings.DB.Transaction(func(tx *gorm.DB) error {
		for _, stream := range streams {
			stream.Job = job
			if err := tx.Save(&stream).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

type language struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func loadLanguagesFromFile(file string) []language {
	// Open json file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read json file
	var languages []language
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&languages)
	if err != nil {
		panic(err)
	}

	return languages
}

func loadSupportedLanguages() []string {
	languages := loadLanguagesFromFile("./misc/supported-languages.json")
	var supportedLanguages []string
	for _, language := range languages {
		supportedLanguages = append(supportedLanguages, language.Code)
	}
	return supportedLanguages
}

func Orchestrator() {
	var wg sync.WaitGroup
	var languages []string = loadSupportedLanguages()

	println("Loading save-stream jobs...")

	// Create a new job
	job := models.Job{
		Name:        "Streams",
		Description: "Streams",
		Status:      "Running",
	}

	// Save job
	if err := settings.DB.Save(&job).Error; err != nil {
		panic(err)
	}

	for _, language := range languages {
		wg.Add(1)
		go func(job models.Job, language string) {
			defer wg.Done()
			LoopStreams(job, []string{language})
		}(job, language)
	}

	wg.Wait()

	println("Finished orchestrating.")

	// Update job
	job.Status = "Finished"
	if err := settings.DB.Save(&job).Error; err != nil {
		panic(err)
	}

}

func LoopStreams(job models.Job, language []string) error {
	var streams int = 0 // Stream count

	params := &twitch.GetStreamListParams{}
	params.Language = language

	// Buffer all pages
	for ok := true; ok; ok = params.After != "" {

		hasRate, when := twitch.GetTwitchService().HasRateLimit()
		if !hasRate {
			println("Rate limit exceeded, waiting...")
			println("Rate limit reset:", when.Format(time.RFC3339))
			duration := when.Sub(time.Now())
			time.Sleep(duration)
		}

		stream, err := twitch.GetTwitchService().GetStreamList(params)
		if err != nil {
			return err
		}

		params.After = stream.Pagination.Cursor

		streams += len(stream.Data)
		UpsertStreams(stream.Data, job)
	}

	println("Streams finished:", streams)

	return nil
}
