package service

import (
	"context"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/imagen/pubsub"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func newImageService(ctx context.Context) domain.ImageService {
	return &imageService{
		pubsubClient: pubsub.NewClient(ctx),
	}
}

type imageService struct {
	pubsubClient pubsub.Client
}

func (s imageService) Generate(ctx context.Context, prompt string) error {
	client, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Generate: %w", err)
	}

	defer client.Close()

	tls, err := client.Translate(ctx, []string{prompt}, language.English, &translate.Options{
		Source: language.Japanese,
	})

	if err != nil {
		return fmt.Errorf("Generate: %w", err)
	}

	if err = s.pubsubClient.PublishGenerateImage(ctx, tls[0].Text); err != nil {
		return fmt.Errorf("Generate: %w", err)
	}

	return nil
}
