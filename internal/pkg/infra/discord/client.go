package discord

import "github.com/go-resty/resty/v2"

type Client struct {
	clientID     string
	clientSecret string
	restClient   *resty.Client
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		restClient:   resty.New(),
	}
}

func (c Client) RegisterCommand() error {
	return nil
}
