package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"imagen/internal/pkg/domain"

	"cloud.google.com/go/pubsub"
)

const (
	topicGenerateImage = "generate-image"
)

type Client struct {
	projectID string
}

func NewClient(projectID string) *Client {
	return &Client{
		projectID: projectID,
	}
}

func (c Client) PublishGenerateImage(ctx context.Context, command domain.ImageGenerateComamnd, extra map[string]interface{}) error {
	client, err := pubsub.NewClient(ctx, c.projectID)
	if err != nil {
		return fmt.Errorf("PublishGenerateImage: projectID=%v: %w", c.projectID, err)
	}

	defer client.Close()

	data, err := json.Marshal(map[string]interface{}{
		"prompt":         command.Prompt,
		"width":          command.Width,
		"height":         command.Height,
		"init_image_url": command.InitImageURL,
		"extra":          extra,
	})

	if err != nil {
		return fmt.Errorf("PublishGenerateImage: %w", err)
	}

	t := client.Topic(topicGenerateImage)
	res := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	if _, err := res.Get(ctx); err != nil {
		return fmt.Errorf("PublishGenerateImage: topic=%v: %w", t.String(), err)
	}

	return nil
}
