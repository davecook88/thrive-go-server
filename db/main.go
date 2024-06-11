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