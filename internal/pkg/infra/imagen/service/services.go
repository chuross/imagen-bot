package service

import (
	"imagen/internal/pkg/domain"
)

type Services struct {
	Image domain.ImageService
}

func NewServices() Services {
	return Services{
		Image: newImageService(),
	}
}
