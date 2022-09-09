package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"imagen/internal/pkg/infra/environment"

	"cloud.google.com/go/pubsub"
)

const (
	topicGenerateImage = "generate-image"
)

type Client interface {
	PublishGenerateImage(ctx context.Context, prompt string) error
}

func NewClient() Client {
	return client{}
}

type client struct {
}

func (c client) PublishGenerateImage(ctx context.Context, prompt string) error {
	env := environment.MustGet(ctx)

	client, err := pubsub.NewClient(ctx, env.GOOGLE_CLOUD_PROJECT_ID)
	if err != nil {
		return fmt.Errorf("PublishGenerateImage: projectID=%v: %w", env.GOOGLE_CLOUD_PROJECT_ID, err)
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
		return fmt.Errorf("PublishGenerateImage: topic=%v: %w", t.String(), err)
	}

	return nil
}
