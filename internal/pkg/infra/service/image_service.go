package service

import (
	"context"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/environment"
	"imagen/internal/pkg/infra/pubsub"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func newImageService() domain.ImageService {
	return &imageService{}
}

type imageService struct {
}

func (s imageService) Generate(ctx context.Context, command domain.ImageGenerateComamnd, extra map[string]interface{}) error {
	translateClient, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Generate: %w", err)
	}
	defer translateClient.Close()

	tls, err := translateClient.Translate(ctx, []string{command.Prompt}, language.English, &translate.Options{
		Source: language.Japanese,
	})

	if err != nil {
		return fmt.Errorf("Generate: %w", err)
	}

	command.Prompt = tls[0].Text

	pubsubClient := pubsub.NewClient(environment.MustGet().GOOGLE_CLOUD_PROJECT_ID)
	if err = pubsubClient.PublishGenerateImage(ctx, command, extra); err != nil {
		return fmt.Errorf("Generate: %w", err)
	}

	return nil
}
