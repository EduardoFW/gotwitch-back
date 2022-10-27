package v1

import (
	"api.gotwitch.tk/models"
	"api.gotwitch.tk/settings"
	"github.com/gin-gonic/gin"
)

func lastSuccessfullyRanJob() (models.Job, error) {
	var job models.Job
	err := settings.DB.Where("status = ?", "Finished").Order("created_at desc").Take(&job).Error
	return job, err
}

func randomStreamFromJob(job models.Job, games []string, languages []string) (models.Stream, error) {
	var stream models.Stream
	query := settings.DB.Model(&stream).Where("job_id = ?", job.ID)
	if len(games) > 0 {
		query = query.Where("game_id in (?)", games)
	}
	if len(languages) > 0 {
		query = query.Where("language in (?)", languages)
	}
	err := query.Order("RANDOM()").Limit(1).Take(&stream).Error
	return stream, err
}

type RandomStreamParams struct {
	Games     []string `url:"game_id[]"`
	Languages []string `url:"language[]"`
}

func GetRandomStream(c *gin.Context) {
	lastJob, err := lastSuccessfullyRanJob()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var params RandomStreamParams
	params.Games = c.QueryArray("game_id[]")
	params.Languages = c.QueryArray("language[]")

	stream, err := randomStreamFromJob(lastJob, params.Games, params.Languages)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{
				"error": "No streams found",
			})
			return
		}

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"stream": stream,
	})
}

func SearchCategories(c *gin.Context) {
	var categories []models.Category
	query := c.Query("query")
	// Search case insensitive
	err := settings.DB.Where("name ILIKE ?", "%"+query+"%").Limit(10).Find(&categories).Error
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"categories": categories,
	})
}
