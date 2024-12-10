package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func BuildEngine() *gin.Engine {
	engine := gin.Default()

	corsConfig := cors.DefaultConfig()

	corsConfig.AllowAllOrigins = true

	engine.Use(cors.New(corsConfig))

	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return engine
}
