package routers

import (
	v1 "api.gotwitch.tk/controllers/v1"
	v2 "api.gotwitch.tk/controllers/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(cors.Default())

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/ping", v1.Ping)

		// RandomStream
		apiV1.GET("/random-stream", v1.GetRandomStream)
		apiV1.GET("/search-category", v1.SearchCategories)
	}

	apiV2 := r.Group("/api/v2")
	{
		// RandomStream
		apiV2.GET("/random-stream", v2.GetRandomStream)
		// apiV1.GET("/search-category", v1.SearchCategories)
	}

	return r
}
