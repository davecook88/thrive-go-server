package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"thrive/server/chatgpt"
	"thrive/server/db"

	"github.com/gin-gonic/gin"
)

func NewChatGPTRequest(messages []chatgpt.Message) *chatgpt.ChatGPTRequest {
	return &chatgpt.ChatGPTRequest{
		Model:    "gpt-4o",
		Messages: messages,
		Stream:   true,
	}
}

func CallChatGPT(c *gin.Context, messages []chatgpt.Message) (*chatgpt.Message, error) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	client := chatgpt.NewChatGPTClient(os.Getenv("OPENAI_API_KEY"))
	jsonData, err := json.Marshal(NewChatGPTRequest(messages))
	if err != nil {
		return nil, errors.New("failed to marshal request")
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.New("failed to create request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to make request")
	}
	defer resp.Body.Close()
	// print raw response as a string
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, errors.New("failed to read response")
	// }

	responseMessage := chatgpt.Message{
		Role:    "assistant",
		Content: "",
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data != "[DONE]" {
				var responseData chatgpt.StreamingResponse
				if err := json.Unmarshal([]byte(data), &responseData); err != nil {
					c.SSEvent("error", gin.H{"error": "Failed to unmarshal response data"})
					return nil, errors.New("failed to unmarshal response data")
				}
				if responseData.Choices[0].Delta.Content != nil {
					responseMessage.Content += *responseData.Choices[0].Delta.Content
					c.SSEvent("message", gin.H{"content": responseMessage.Content})
				}
			}
		}
	}

	return &responseMessage, nil

}

func ChatGPTHandler(c *gin.Context) {
	var request UserMessage

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dbClient, err := db.NewClient(c, "thrive-chat")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var existingMessages *[]chatgpt.Message

	if request.ChatID != nil {
		existingMessages, err = dbClient.GetChat(c, *request.ChatID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if existingMessages == nil {
		existingMessages = &[]chatgpt.Message{}
	}

	messages := append(*existingMessages, chatgpt.Message{Role: chatgpt.UserRole, Content: request.Message})

	chatGPTResponseMessage, err := CallChatGPT(c, messages)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	messages = append(messages, *chatGPTResponseMessage)

	if err := dbClient.CreateChat(c, messages); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Marshal the chatGPTResponse struct back to JSON
	jsonResponse, err := json.Marshal(messages)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal ChatGPT response"})
		return
	}

	// Print the full JSON response
	fmt.Println(string(jsonResponse))

	c.JSON(200, messages)

}
