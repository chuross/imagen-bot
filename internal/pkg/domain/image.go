package domain

import "context"

type ImageGenerateComamnd struct {
	Prompt          string
	NegativePrompts []string
	RawPrompt       string
	Width           int
	Height          int
	InitImageURL    *string
	MaskImageURL    *string
	Strength        float64
	Number          int
}

type ImageService interface {
	Generate(ctx context.Context, command ImageGenerateComamnd, extra map[string]interface{}) error
}
