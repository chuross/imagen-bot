package domain

import "context"

type ImageService interface {
	Generate(ctx context.Context, prompt string, extra map[string]interface{}) error
}
