package domain

import "context"

type ImageGenerateComamnd struct {
	Prompt       string
	Width        int
	Height       int
	InitImageURL *string
	MaskImageURL *string
	Strength     float64
}

type ImageService interface {
	Generate(ctx context.Context, command ImageGenerateComamnd, extra map[string]interface{}) error
	Upscale(ctx context.Context, imageURL string, extra map[string]interface{}) error
}
