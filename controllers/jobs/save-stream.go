package jobs

import (
	"encoding/json"
	"os"
	"time"

	"api.gotwitch.tk/models"
	"api.gotwitch.tk/services/twitch"
	"api.gotwitch.tk/settings"
	"github.com/gammazero/workerpool"
	"github.com/getsentry/sentry-go"
	"gorm.io/gorm/clause"
)

var streamWorkerPoolSize = 3

func UpsertStreams(streams []models.Stream, job models.Job) error {
	if len(streams) == 0 {
		return nil
	}
	// Set the job for the streams
	for i := range streams {
		streams[i].Job = job
	}

	return settings.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&streams).Error

}

func DeleteStreamsWithJobIDLessThan(id uint) error {
	return settings.DB.Where("job_id < ?", id).Unscoped().Delete(&models.Stream{}).Error
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

	runWorkerPool(languages, job)

	println("Finished orchestrating.")

	// Update job
	job.Status = "Finished"
	if err := settings.DB.Save(&job).Error; err != nil {
		panic(err)
	}

	// Delete streams with job_id < job.id
	if err := DeleteStreamsWithJobIDLessThan(job.ID); err != nil {
		panic(err)
	}

}

func runWorkerPool(languages []string, job models.Job) {
	wp := workerpool.New(streamWorkerPoolSize)

	for _, language := range languages {

		language := language
		job := job

		wp.Submit(func() {
			LoopStreams(job, []string{language})
		})
	}

	wp.StopWait()
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
			sentry.CaptureException(err)
			break
		}

		params.After = stream.Pagination.Cursor

		streams += len(stream.Data)
		UpsertStreams(stream.Data, job)
	}

	println("Streams finished:", streams)

	return nil
}
