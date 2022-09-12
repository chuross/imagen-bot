package domain

import "context"

type ImageGenerateComamnd struct {
	Prompt          string
	Width           int
	Height          int
	InitImageBase64 *string
}

type ImageService interface {
	Generate(ctx context.Context, command ImageGenerateComamnd, extra map[string]interface{}) error
}
