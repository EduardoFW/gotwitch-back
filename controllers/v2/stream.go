package v1

import (
	"api.gotwitch.tk/models"
	"api.gotwitch.tk/settings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func lastSuccessfullyRanJob() (models.Job, error) {
	var job models.Job
	err := settings.DB.Where("status = ?", "Finished").Order("created_at desc").Take(&job).Error
	return job, err
}

func randomStreamFromJob(job models.Job, params RandomStreamParams) (models.Stream, error) {
	var stream models.Stream
	query := settings.DB.Model(&stream).Where("job_id = ?", job.ID)
	query = filterQueryFromParams(query, params)
	err := query.Order("RANDOM()").Limit(1).Take(&stream).Error
	return stream, err
}

func filterQueryFromParams(query *gorm.DB, params RandomStreamParams) *gorm.DB {
	if len(params.Games) > 0 {
		query = query.Where("game_id in (?)", params.Games)
	}
	if len(params.Languages) > 0 {
		query = query.Where("language in (?)", params.Languages)
	}
	return query
}

func buildQueryFromParams(params RandomStreamParams) *gorm.DB {
	query := settings.DB.Model(&models.Stream{})
	query = filterQueryFromParams(query, params)
	return query
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

	params := getRandomStreamParamsFromGinContext(c)

	stream, err := randomStreamFromJob(lastJob, params)
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

func getRandomStreamParamsFromGinContext(c *gin.Context) RandomStreamParams {
	var params RandomStreamParams
	params.Games = c.QueryArray("game_id[]")
	params.Languages = c.QueryArray("language[]")
	return params
}

func GetStreamCount(c *gin.Context) {
	var count int64
	lastJob, err := lastSuccessfullyRanJob()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	params := getRandomStreamParamsFromGinContext(c)

	query := buildQueryFromParams(params)

	err = query.Where("job_id = ?", lastJob.ID).Count(&count).Error
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"count": count,
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
