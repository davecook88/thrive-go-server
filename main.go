package main

import (
	"thrive/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
		panic(err)
	}
    r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/chatgpt", handlers.ChatGPTHandler)

	r.Run() // listen and serve on 0.0.0.0:8080
}