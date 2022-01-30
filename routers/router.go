package routers

import (
	v1 "github.com/EduardoWeber/twitch-go/controllers/v1"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/ping", v1.Ping)

		// RandomStream
		apiV1.GET("/random-stream", v1.GetRandomStream)
	}

	return r
}
