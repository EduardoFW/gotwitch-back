package jobs

import (
	"time"

	"api.gotwitch.tk/models"
	"api.gotwitch.tk/services/twitch"
	"api.gotwitch.tk/settings"
	"gorm.io/gorm"
)

func UpsertStreams(streams []models.Stream) error {
	return settings.DB.Transaction(func(tx *gorm.DB) error {
		for _, stream := range streams {

			if err := tx.Save(&stream).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func LoopStreams() error {
	var streams int = 0 // Stream count

	params := &twitch.GetStreamListParams{}

	// Buffer all pages
	for ok := true; ok; ok = params.After != "" {

		// params.Language = []string{"el"}

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
		UpsertStreams(stream.Data)
	}

	println("Streams finished:", streams)

	return nil
}
