package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

const (
	topicGenerateImage = "generate-image"
)

type Client interface {
	PublishGenerateImage(ctx context.Context, prompt string) error
}

func NewClient(ctx context.Context) Client {
	return &client{}
}

type client struct {
}

func (c client) PublishGenerateImage(ctx context.Context, prompt string) error {
	client, err := pubsub.NewClient(ctx, "project-id")
	if err != nil {
		return fmt.Errorf("PublishGenerateImage: %w", err)
	}

	defer client.Close()

	data, err := json.Marshal(map[string]string{
		"prompt": prompt,
	})

	if err != nil {
		return fmt.Errorf("PublishGenerateImage: %w", err)
	}

	t := client.Topic(topicGenerateImage)
	res := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	if _, err := res.Get(ctx); err != nil {
		return fmt.Errorf("PublishGenerateImage: %w", err)
	}

	return nil
}
