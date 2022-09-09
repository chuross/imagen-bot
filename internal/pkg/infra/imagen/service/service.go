package service

import (
	"context"
	"imagen/internal/pkg/domain"
)

type Services struct {
	Image domain.ImageService
}

func NewServices(ctx context.Context) Services {
	return Services{
		Image: newImageService(ctx),
	}
}
