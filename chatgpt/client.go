package chatgpt

import (
	"net/http"
)

type ChatGPTClient struct {
	ApiKey string
	*http.Client
}

func NewChatGPTClient(apiKey string) *ChatGPTClient {
	return &ChatGPTClient{
		ApiKey: apiKey,
		Client: &http.Client{},
	}
}

func (c *ChatGPTClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	return c.Client.Do(req)
}