package webhook

import (
	"context"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/imagen/service"
	"log"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/mattn/go-shellwords"
	"github.com/samber/lo"
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
			sendingTargetID := event.Source.UserID

			args, err := shellwords.Parse(event.Message.(*linebot.TextMessage).Text)
			if err != nil {
				return fmt.Errorf("GenerateByLine: %w", err)
			}

			text := args[0]
			width, height, err := u.resolveSize(args)
			if err != nil {
				return fmt.Errorf("GenerateByLine: %w", err)
			}

			command := domain.ImageGenerateComamnd{
				Prompt: text,
				Width:  width,
				Height: height,
			}

			if err := u.imageService.Generate(ctx, command, map[string]interface{}{
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

func (u ImageUseCase) resolveSize(args []string) (int, int, error) {
	sizeOpt, found := lo.Find(args, func(arg string) bool {
		return strings.HasPrefix(arg, "-s")
	})

	if !found {
		return 0, 0, nil
	}

	size := strings.Split(strings.TrimSpace(sizeOpt[2:]), "x")

	width, err := strconv.Atoi(size[0])
	if err != nil {
		return 0, 0, fmt.Errorf("resolveSize: %w", err)
	}

	height, err := strconv.Atoi(size[0])
	if err != nil {
		return 0, 0, fmt.Errorf("resolveSize: %w", err)
	}

	return width, height, nil
}
