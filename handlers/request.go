package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"thrive/server/chatgpt"
	"thrive/server/db"

	"github.com/gin-gonic/gin"
)




func ChatGPTHandler(c *gin.Context) {
	var request chatgpt.ChatGPTRequest


	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(request); if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal ChatGPT request"})
		return
	}

	client := chatgpt.NewChatGPTClient(os.Getenv("OPENAI_API_KEY"))
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create request"})
        return
    }
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to make request"})
		return
	}
	defer resp.Body.Close()

	dbClient, err := db.NewClient(c, "thrive-chat"); if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := dbClient.CreateChat(c, request.Messages); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

    var chatGPTResponse chatgpt.ChatGPTResponse
    if err := json.NewDecoder(resp.Body).Decode(&chatGPTResponse); err != nil {
        c.JSON(500, gin.H{"error": "Failed to decode ChatGPT response"})
        return
    }



    // Marshal the chatGPTResponse struct back to JSON
    jsonResponse, err := json.Marshal(chatGPTResponse)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to marshal ChatGPT response"})
        return
    }

    // Print the full JSON response
    fmt.Println(string(jsonResponse))

    c.JSON(200, chatGPTResponse)

}