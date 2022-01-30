package v1

import (
	"math/rand"

	"github.com/EduardoWeber/twitch-go/models"
	"github.com/EduardoWeber/twitch-go/services"
	"github.com/EduardoWeber/twitch-go/settings"
	"github.com/gin-gonic/gin"
)

func GetRandomStream(c *gin.Context) {
	maxPageSize := 10

	token, err := services.GetTwitchToken(settings.ServerSettings.TwitchClientID, settings.ServerSettings.TwitchClientSecret)
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

	for i := 0; i < pages; i++ {

		stream, err := services.GetStreamList(token, settings.ServerSettings.TwitchClientID, after)
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
