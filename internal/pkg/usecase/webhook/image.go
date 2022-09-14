package webhook

import (
	"context"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/environment"
	"imagen/internal/pkg/infra/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type imageGenerateOption struct {
	Text   string
	Width  int
	Height int
}

type ImageUseCase struct {
	imageService domain.ImageService
}

func newImageUseCase(services *service.Services) *ImageUseCase {
	return &ImageUseCase{
		imageService: services.Image,
	}
}

func (u ImageUseCase) GenerateByDiscordMessageCommand(ctx context.Context, interact *discordgo.Interaction) error {
	data := interact.ApplicationCommandData()

	if !strings.HasPrefix(data.Name, "imagen") {
		return fmt.Errorf("GenerateByDiscord: unexpected command: %v", interact.ApplicationCommandData().Name)
	}

	message, ok := data.Resolved.Messages[data.TargetID]
	if !ok {
		return nil
	}

	var initImageURL *string
	var maskImageURL *string

	if messageRef := message.MessageReference; messageRef != nil {
		discordSes, err := discordgo.New(fmt.Sprintf("Bot %s", environment.MustGet().DISCORD.BOT_TOKEN))
		if err != nil {
			return fmt.Errorf("GenerateByDiscord: %w", err)
		}

		referencedMes, err := discordSes.ChannelMessage(messageRef.ChannelID, messageRef.MessageID)
		if err != nil {
			return fmt.Errorf("GenerateByDiscord: %w", err)
		}

		if len(referencedMes.Attachments) > 0 {
			initImageURL = &referencedMes.Attachments[0].URL
		}

		if len(message.Attachments) > 0 {
			maskImageURL = &message.Attachments[0].URL
		}
	} else {
		if len(message.Attachments) > 0 {
			initImageURL = &message.Attachments[0].URL
		}
	}

	command := domain.ImageGenerateComamnd{
		Prompt:       message.Content,
		Width:        0,
		Height:       0,
		InitImageURL: initImageURL,
		MaskImageURL: maskImageURL,
	}

	if err := u.imageService.Generate(ctx, command, map[string]interface{}{
		"via":               "discord",
		"user_id":           interact.Member.User.ID,
		"interaction_id":    interact.ID,
		"interaction_token": interact.Token,
		"message_id":        message.ID,
		"message_url":       fmt.Sprintf("https://discord.com/channels/%s/%s/%s", interact.GuildID, interact.ChannelID, message.ID),
	}); err != nil {
		return fmt.Errorf("GenerateByDiscord: %w", err)
	}

	return nil
}

func (u ImageUseCase) GenerateByDiscordMessageComponent(ctx context.Context, interact *discordgo.Interaction) error {
	return nil
}
