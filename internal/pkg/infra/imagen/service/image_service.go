package service

import (
	"context"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/imagen/pubsub"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func newImageService() domain.ImageService {
	return imageService{}
}

type imageService struct {
}

func (s imageService) Generate(ctx context.Context, prompt string) error {
	translateClient, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Generate: %w", err)
	}
	defer translateClient.Close()

	tls, err := translateClient.Translate(ctx, []string{prompt}, language.English, &translate.Options{
		Source: language.Japanese,
	})

	if err != nil {
		return fmt.Errorf("Generate: %w", err)
	}

	pubsubClient := pubsub.NewClient()
	if err = pubsubClient.PublishGenerateImage(ctx, tls[0].Text); err != nil {
		return fmt.Errorf("Generate: %w", err)
	}

	return nil
}
