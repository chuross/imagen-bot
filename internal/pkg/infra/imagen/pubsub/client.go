package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

const (
	topicGenerateImage = "generate_image"
)

type Client interface {
}

func NewClient(ctx context.Context) Client {
	c, err := pubsub.NewClient(ctx, "project-id")
	if err != nil {
		panic(fmt.Errorf("pubsub initialize failed: %w", err).Error())
	}
	return &client{
		client: c,
	}
}

type client struct {
	client *pubsub.Client
}

func (c client) PublishGenerateImage(ctx context.Context, prompt string) error {
	data, err := json.Marshal(map[string]string{
		"prompt": prompt,
	})

	if err != nil {
		return fmt.Errorf("publish generate image failed: %w", err)
	}

	t := c.client.Topic(topicGenerateImage)
	res := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	if _, err := res.Get(ctx); err != nil {
		return fmt.Errorf("publish generate image failed: %w", err)
	}

	return nil
}
