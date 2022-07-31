package v1

import (
	"math/rand"
	"strconv"

	"api.gotwitch.tk/models"
	"api.gotwitch.tk/services/twitch"
	"github.com/gin-gonic/gin"
)

func GetRandomStream(c *gin.Context) {
	maxPageSize := 10

	// Random number between 1 and pageSize
	pages := rand.Intn(maxPageSize) + 1

	// Streams list
	var streams []models.Stream
	var after string

	params := &twitch.GetStreamListParams{}

	// Buffer all pages
	for i := 0; i < pages; i++ {

		params.After = after
		params.Language = c.QueryArray("language[]")
		params.GameId = c.QueryArray("game_id[]")
		stream, err := twitch.GetTwitchService().GetStreamList(params)
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
	if len(streams) == 0 {
		c.JSON(500, gin.H{
			"error": "No streams found",
		})
		return
	}

	index := rand.Intn(len(streams))

	c.JSON(200, gin.H{
		"data": streams[index],
	})
}

func SearchCategories(c *gin.Context) {
	first, _ := strconv.Atoi(c.DefaultQuery("first", "100"))
	params := &twitch.SearchCategoriesParams{
		Query: c.Query("query"),
		First: first,
		After: c.Query("after"),
	}

	resp, err := twitch.GetTwitchService().SearchCategories(params)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, resp)
}
