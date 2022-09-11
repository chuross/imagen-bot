package webhook

import (
	"context"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/imagen/service"
	"log"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/mattn/go-shellwords"
)

type ImageGenerateOption struct {
	Size string `short:"s" long:"size" description:"widthxheight (ex)256x256"`
}

type ImageUseCase struct {
	imageService domain.ImageService
}

func newImageUseCase(services service.Services) *ImageUseCase {
	return &ImageUseCase{
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

			text := event.Message.(*linebot.TextMessage).Text
			args, err := shellwords.Parse(text)
			if err != nil {
				return fmt.Errorf("GenerateByLine: %w", err)
			}

			var opt ImageGenerateOption
			if _, err := flags.ParseArgs(&opt, args); err != nil {
				return fmt.Errorf("GenerateByLine: %w", err)
			}

			text = strings.Split(text, "-")[0]

			width, height, err := u.resolveSize(opt)
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

func (u ImageUseCase) resolveSize(opt ImageGenerateOption) (int, int, error) {
	if opt.Size == "" {
		return 0, 0, nil
	}

	size := strings.Split(strings.TrimSpace(opt.Size), "x")

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
