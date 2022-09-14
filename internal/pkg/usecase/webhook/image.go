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
		return fmt.Errorf("GenerateByDiscordMessageCommand: unexpected command: %v", interact.ApplicationCommandData().Name)
	}

	message, ok := data.Resolved.Messages[data.TargetID]
	if !ok {
		return nil
	}

	if err := u.generate(ctx, interact.GuildID, interact.ChannelID, message.Author.ID, interact.Token, data.Name, message); err != nil {
		return fmt.Errorf("GenerateByDiscordMessageCommand: %w", err)
	}

	return nil
}

func (u ImageUseCase) GenerateByDiscordMessageComponent(ctx context.Context, interact *discordgo.Interaction) error {
	data := interact.MessageComponentData()
	customID := data.CustomID

	params := strings.Split(customID, "##")
	if len(params) != 2 {
		return fmt.Errorf("GenerateByDiscordMessageComponent: invalid custom id: id=%v", customID)
	}

	messageID := params[1]

	session, err := discordgo.New(fmt.Sprintf("Bot %s", environment.MustGet().DISCORD.BOT_TOKEN))
	if err != nil {
		return nil
	}

	message, err := session.ChannelMessage(interact.ChannelID, messageID)
	if err != nil {
		return fmt.Errorf("GenerateByDiscordMessageComponent: %w", err)
	}

	if err := u.generate(ctx, interact.GuildID, interact.ChannelID, message.Author.ID, interact.Token, message.ReferencedMessage.Interaction.Name, message); err != nil {
		return fmt.Errorf("GenerateByDiscordMessageComponent: %w", err)
	}

	return nil
}

func (u ImageUseCase) generate(ctx context.Context, guildID, channelID, userID, interactionToken, commandName string, message *discordgo.Message) error {
	var initImageURL *string

	if len(message.Attachments) > 0 {
		initImageURL = &message.Attachments[0].URL
	}

	if messageRef := message.MessageReference; messageRef != nil {
		discordSes, err := discordgo.New(fmt.Sprintf("Bot %s", environment.MustGet().DISCORD.BOT_TOKEN))
		if err != nil {
			return fmt.Errorf("generate: %w", err)
		}

		referencedMes, err := discordSes.ChannelMessage(messageRef.ChannelID, messageRef.MessageID)
		if err != nil {
			return fmt.Errorf("generate: %w", err)
		}

		if len(referencedMes.Attachments) > 0 {
			initImageURL = &referencedMes.Attachments[0].URL
		}
	}

	command := domain.ImageGenerateComamnd{
		Prompt:       message.Content,
		Width:        0,
		Height:       0,
		InitImageURL: initImageURL,
		MaskImageURL: nil,
	}

	if err := u.imageService.Generate(ctx, command, map[string]interface{}{
		"via":               "discord",
		"user_id":           userID,
		"interaction_token": interactionToken,
		"command_name":      commandName,
		"message_id":        message.ID,
		"message_url":       fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, message.ID),
	}); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	return nil
}
