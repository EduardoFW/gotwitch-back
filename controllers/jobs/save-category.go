package jobs

import (
	"time"

	"api.gotwitch.tk/models"
	"api.gotwitch.tk/services/twitch"
	"api.gotwitch.tk/settings"
	"github.com/gammazero/workerpool"
	"github.com/getsentry/sentry-go"
	"gorm.io/gorm/clause"
)

var categoryWorkerPoolSize = 2

func CategoryOrchestrator() {
	println("Starting category orchestrator")

	// Create a new job
	job := models.Job{
		Name:        "Categories",
		Description: "Fetch all categories from Twitch",
		Status:      "Running",
	}
	if err := settings.DB.Create(&job).Error; err != nil {
		sentry.CaptureException(err)
	}

	alphabet := generateAlphabet()

	allCategories := aggregateAllCategoriesWorkerPool(alphabet)

	// Upsert all categories
	err := UpsertCategories(allCategories, job)
	if err != nil {
		sentry.CaptureException(err)
	}

	job.Status = "Finished"
	if err := settings.DB.Save(&job).Error; err != nil {
		sentry.CaptureException(err)
	}

	println("Finished category orchestrator")
}

func UpsertCategories(allCategories []models.Category, job models.Job) error {
	if len(allCategories) == 0 {
		return nil
	}

	// Set the job for the categories
	for i := range allCategories {
		allCategories[i].Job = job
	}

	if err := settings.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(allCategories, 100).Error; err != nil {
		return err
	}

	return nil
}

func aggregateAllCategoriesWorkerPool(letters []string) []models.Category {
	var allCategories []models.Category
	workerChannel := make(chan []models.Category)
	wp := workerpool.New(categoryWorkerPoolSize)

	for _, letter := range letters {
		letter := letter
		wp.Submit(func() {
			categories := getAllCategoriesContaining(letter)
			workerChannel <- categories
		})
	}

	for i := 0; i < len(letters); i++ {
		categories := <-workerChannel
		allCategories = appendUniqueCategories(allCategories, categories)
	}

	wp.StopWait()

	return allCategories
}

func appendUniqueCategories(categoriesDest []models.Category, categoriesToAdd []models.Category) []models.Category {
	for _, categoryToAdd := range categoriesToAdd {
		if !containsCategory(categoriesDest, categoryToAdd) {
			categoriesDest = append(categoriesDest, categoryToAdd)
		}
	}

	return categoriesDest
}

func containsCategory(allCategories []models.Category, category models.Category) bool {
	for _, c := range allCategories {
		if c.Id == category.Id {
			return true
		}
	}

	return false
}

func generateAlphabet() []string {
	var alphabet []string
	for i := 'a'; i <= 'z'; i++ {
		alphabet = append(alphabet, string(i))
	}

	return alphabet
}

func getAllCategoriesContaining(query string) []models.Category {
	println("Fetching categories for letter: " + query)
	var allCategories []models.Category
	params := &twitch.SearchCategoriesParams{
		First: 100,
		Query: query,
	}

	for ok := true; ok; ok = params.After != "" {

		hasRate, when := twitch.GetTwitchService().HasRateLimit()
		if !hasRate {
			println("Rate limit exceeded, waiting...")
			println("Rate limit reset:", when.Format(time.RFC3339))
			duration := when.Sub(time.Now())
			time.Sleep(duration)
		}

		categories, err := twitch.GetTwitchService().SearchCategories(params)
		if err != nil {
			sentry.CaptureException(err)
		}

		allCategories = append(allCategories, categories.Data...)

		params.After = categories.Pagination.Cursor
	}

	return allCategories
}
