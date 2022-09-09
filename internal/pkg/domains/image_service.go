package domains

import "context"

type ImageService interface {
}

func NewImageService(ctx context.Context) ImageService {
	return &imageService{}
}

type imageService struct {
}
