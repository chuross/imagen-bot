package webhook

import "imagen/internal/pkg/infra/imagen/service"

type UseCases struct {
	Image *ImageUseCase
}

func NewWebhookUseCases(services *service.Services) *UseCases {
	return &UseCases{
		Image: newImageUseCase(services),
	}
}
