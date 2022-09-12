package webhook

import (
	"context"
	"encoding/base64"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/imagen/service"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/jessevdk/go-flags"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/mattn/go-shellwords"
)

type imageGenerateOption struct {
	Text   string
	Width  int
	Height int
}

type ImageUseCase struct {
	imageService domain.ImageService
	client       *resty.Client
}

func newImageUseCase(services *service.Services) *ImageUseCase {
	return &ImageUseCase{
		imageService: services.Image,
		client:       resty.New(),
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

			opt, err := resolveImageGenerateOption(event.Message.(*linebot.TextMessage).Text)
			if err != nil {
				return fmt.Errorf("GenerateByLine: %w", err)
			}

			command := domain.ImageGenerateComamnd{
				Prompt: opt.Text,
				Width:  opt.Width,
				Height: opt.Height,
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

func (u ImageUseCase) GenerateByDiscord(ctx context.Context, interact *discordgo.Interaction) error {
	data := interact.ApplicationCommandData()

	if data.Name != "imagen" {
		return fmt.Errorf("GenerateByDiscord: unexpected command: %v", interact.ApplicationCommandData().Name)
	}

	message, ok := data.Resolved.Messages[data.TargetID]
	if !ok {
		return nil
	}

	opt, err := resolveImageGenerateOption(message.Content)
	if err != nil {
		return fmt.Errorf("GenerateByDiscord: %w", err)
	}

	var initImageBase64 *string
	if len(message.Attachments) > 0 {
		attachment := message.Attachments[0]

		if strings.HasPrefix(attachment.ContentType, "image/") {
			if res, err := u.client.R().Get(attachment.URL); err != nil {
				return fmt.Errorf("GenerateByDiscord: %w", err)
			} else {
				d := base64.StdEncoding.EncodeToString(res.Body())
				initImageBase64 = &d
			}
		}
	}

	command := domain.ImageGenerateComamnd{
		Prompt:          opt.Text,
		Width:           opt.Width,
		Height:          opt.Height,
		InitImageBase64: initImageBase64,
	}

	if err := u.imageService.Generate(ctx, command, map[string]interface{}{
		"via":               "discord",
		"interaction_id":    interact.ID,
		"interaction_token": interact.Token,
	}); err != nil {
		return fmt.Errorf("GenerateByDiscord: %w", err)
	}

	return nil
}

func resolveImageGenerateOption(text string) (*imageGenerateOption, error) {
	var o struct {
		Size string `short:"s" long:"size" description:"widthxheight (ex)256x256"`
	}

	args, err := shellwords.Parse(text)
	if err != nil {
		return nil, fmt.Errorf("resolveImageGenerateOption: %w", err)
	}

	if _, err := flags.ParseArgs(&o, args); err != nil {
		return nil, fmt.Errorf("resolveImageGenerateOption: %w", err)
	}

	width, height, err := resolveImageSize(o.Size)
	if err != nil {
		return nil, fmt.Errorf("resolveImageGenerateOption: %w", err)
	}

	return &imageGenerateOption{
		Text:   strings.Split(text, "-")[0],
		Width:  width,
		Height: height,
	}, nil
}

func resolveImageSize(sizeStr string) (int, int, error) {
	if sizeStr == "" {
		return 0, 0, nil
	}

	size := strings.Split(strings.TrimSpace(sizeStr), "x")

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
