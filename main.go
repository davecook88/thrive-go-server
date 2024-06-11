package main

import (
	"net/http"
	"thrive/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.Use(corsMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/chatgpt", handlers.ChatGPTHandler)

	r.Run() // listen and serve on 0.0.0.0:8080
}
