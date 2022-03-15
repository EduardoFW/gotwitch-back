package v1

import (
	"math/rand"
	"strconv"

	"api.gotwitch.tk/controllers"
	"api.gotwitch.tk/models"
	"api.gotwitch.tk/services"
	"api.gotwitch.tk/settings"
	"github.com/gin-gonic/gin"
)

func GetRandomStream(c *gin.Context) {
	maxPageSize := 10

	tokenString, err := controllers.GetTokenManager().GetToken()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Random number between 1 and pageSize
	pages := rand.Intn(maxPageSize) + 1

	// Streams list
	var streams []models.Stream
	var after string

	params := &services.GetStreamListParams{}

	// Buffer all pages
	for i := 0; i < pages; i++ {

		params.After = after
		params.Language = c.QueryArray("language")
		params.GameId = c.QueryArray("game_id")
		stream, err := services.GetStreamList(tokenString, settings.ServerSettings.TwitchClientID, params)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		after = stream.Pagination.Cursor
		streams = append(streams, stream.Data...)
	}

	// Random number between 0 and len(streams)
	index := rand.Intn(len(streams))

	c.JSON(200, gin.H{
		"data": streams[index],
	})
}

func SearchCategories(c *gin.Context) {
	tokenString, err := controllers.GetTokenManager().GetToken()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	first, _ := strconv.Atoi(c.DefaultQuery("first", "100"))
	params := &services.SearchCategoriesParams{
		Query: c.Query("query"),
		First: first,
		After: c.Query("after"),
	}

	resp, err := services.SearchCategories(tokenString, settings.ServerSettings.TwitchClientID, params)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, resp)
}
