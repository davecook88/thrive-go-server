package db

import (
	"context"
	"thrive/server/chatgpt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Client struct {
	*firestore.Client
}


func NewClient(ctx context.Context, projectID string) (*Client, error) {
	opt := option.WithCredentialsFile("creds.json")

    app, err := firebase.NewApp(ctx, nil, opt)
    if err != nil {
        return nil, err
    }

    client, err := app.Firestore(ctx)
    if err != nil {
        return nil, err
    }

    return &Client{Client: client}, nil

}


func (c *Client) CreateChat(ctx context.Context, messages []chatgpt.Message ) error {
	chatDoc := map[string]interface{}{
        "messages": messages,
        // Add any other fields you want to store in the document
    }
	_, _, err := c.Collection("thrive-chats").Add(ctx, chatDoc)
	return err
}

func (c *Client) UpdateChat(ctx context.Context, chatId string,  messages []chatgpt.Message) error {
	chatDoc := map[string]interface{}{
		"messages": messages,
	}
	_, err := c.Collection("thrive-chats").Doc(chatId).Set(ctx, chatDoc, firestore.MergeAll)
	return err
}

func (c *Client) GetChat(ctx context.Context, chatId string) (*[]chatgpt.Message, error) {
	doc, err := c.Collection("thrive-chats").Doc(chatId).Get(ctx)
	if err != nil {
		return nil, err
	}
	var chat struct {
		Messages []chatgpt.Message `firestore:"messages"`
	}
	if err := doc.DataTo(&chat); err != nil {
		return nil, err
	}
	return &chat.Messages, nil
}