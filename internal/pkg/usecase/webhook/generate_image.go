package webhook

import (
	"context"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/imagen/service"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type ImageUseCase struct {
	imageService domain.ImageService
}

func newImageUseCase(services service.Services) ImageUseCase {
	return ImageUseCase{
		imageService: services.Image,
	}
}

func (u ImageUseCase) GenerateByLine(ctx context.Context, events []*linebot.Event) error {
	for _, event := range events {
		if event.Type != linebot.EventTypeMessage {
			log.Printf("GenerateByLine: unexpected event type: type=%v", event.Type)
			continue
		}

		switch event.Message.(type) {
		case *linebot.TextMessage:
			text := event.Message.(*linebot.TextMessage).Text
			sendingTargetID := event.Source.UserID

			if err := u.imageService.Generate(ctx, text, map[string]interface{}{
				"via":               "line-bot",
				"sending_target_id": sendingTargetID,
				"reply_token":       event.ReplyToken,
			}); err != nil {
				return fmt.Errorf("GenerateByLine: %w", err)
			}
		default:
			fmt.Printf("GenerateByLine: unexpected message type: type=%v", event.Message)
		}
	}

	return nil
}
