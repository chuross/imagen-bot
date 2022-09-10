package webhook

import (
	"context"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/imagen/service"
)

type ImageUseCase struct {
	imageService domain.ImageService
}

func newImageUseCase(services service.Services) ImageUseCase {
	return ImageUseCase{
		imageService: services.Image,
	}
}

func (u ImageUseCase) Generate(ctx context.Context, text, sendingTargetID, replyToken string) error {
	return u.imageService.Generate(ctx, text, map[string]interface{}{
		"via":               "line-bot",
		"sending_target_id": sendingTargetID,
		"reply_token":       replyToken,
	})
}
