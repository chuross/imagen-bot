package domain

import "context"

type ImageGenerateComamnd struct {
	Prompt string
	Width  int
	Height int
}

type ImageService interface {
	Generate(ctx context.Context, command ImageGenerateComamnd, extra map[string]interface{}) error
}
