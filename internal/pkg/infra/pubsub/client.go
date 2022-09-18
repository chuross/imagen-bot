package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"imagen/internal/pkg/domain"
	"log"

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
	data, err := json.Marshal(map[string]interface{}{
		"event_name": "generate",
		"params": map[string]interface{}{
			"prompt":         command.Prompt,
			"width":          command.Width,
			"height":         command.Height,
			"strength":       command.Strength,
			"init_image_url": command.InitImageURL,
			"mask_image_url": command.MaskImageURL,
			"extra":          extra,
		},
	})

	if err != nil {
		return fmt.Errorf("PublishGenerateImage: %w", err)
	}

	if err := c.publish(ctx, data); err != nil {
		return fmt.Errorf("PublishGenerateImage: %w", err)
	}

	return nil
}

func (c Client) PublishUpscaleImage(ctx context.Context, imageURL string, extra map[string]interface{}) error {
	data, err := json.Marshal(map[string]interface{}{
		"event_name": "upscaling",
		"params": map[string]interface{}{
			"image_url": imageURL,
			"extra":     extra,
		},
	})

	if err != nil {
		return fmt.Errorf("PublishUpscaleImage: %w", err)
	}

	if err := c.publish(ctx, data); err != nil {
		return fmt.Errorf("PublishUpscaleImage: %w", err)
	}

	return nil
}

func (c Client) publish(ctx context.Context, data []byte) error {
	client, err := pubsub.NewClient(ctx, c.projectID)
	if err != nil {
		return fmt.Errorf("publish: projectID=%v: %w", c.projectID, err)
	}

	defer client.Close()

	t := client.Topic(topicGenerateImage)
	res := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	if _, err := res.Get(ctx); err != nil {
		return fmt.Errorf("publish: topic=%v: %w", t.String(), err)
	}

	log.Println("pubsub_client: publish successful")

	return nil
}
