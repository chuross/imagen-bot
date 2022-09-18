package webhook

import (
	"context"
	"errors"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/commandline"
	"imagen/internal/pkg/infra/discord"
	"imagen/internal/pkg/infra/environment"
	"imagen/internal/pkg/infra/service"
	"strconv"
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

func (u ImageUseCase) UpscaleByMessageCommand(ctx context.Context, interact *discordgo.Interaction) error {
	data := interact.ApplicationCommandData()

	if data.Name != discord.CommandImagenUpscaling.Name {
		return fmt.Errorf("UpscaleByMessageCommand: unexpected command: %v", interact.ApplicationCommandData().Name)
	}

	message, ok := data.Resolved.Messages[data.TargetID]
	if !ok {
		return nil
	}

	if len(message.Attachments) == 0 {
		return errors.New("UpscaleByMessageCommand: no attachment")
	}

	imageURL := message.Attachments[0].URL
	u.imageService.Upscale(ctx, imageURL, imagenExtra(interact.Token, interact.GuildID, interact.ChannelID, message.Author.ID, message.ID))

	return nil
}

func (u ImageUseCase) GenerateByMessageCommand(ctx context.Context, interact *discordgo.Interaction) error {
	data := interact.ApplicationCommandData()

	if data.Name != discord.CommandImagen.Name {
		return fmt.Errorf("GenerateByMessageCommand: unexpected command: %v", interact.ApplicationCommandData().Name)
	}

	message, ok := data.Resolved.Messages[data.TargetID]
	if !ok {
		return nil
	}

	if err := u.generate(ctx, interact.GuildID, interact.ChannelID, message.Author.ID, interact.Token, message); err != nil {
		return fmt.Errorf("GenerateByMessageCommand: %w", err)
	}

	return nil
}

func (u ImageUseCase) GenerateByMessageComponent(ctx context.Context, interact *discordgo.Interaction) error {
	data := interact.MessageComponentData()
	customID := data.CustomID

	params := strings.Split(customID, "##")
	if len(params) != 2 {
		return fmt.Errorf("GenerateByMessageComponent: invalid custom id: id=%v", customID)
	}

	messageID := params[1]

	session, err := discordgo.New(fmt.Sprintf("Bot %s", environment.MustGet().DISCORD.BOT_TOKEN))
	if err != nil {
		return nil
	}

	message, err := session.ChannelMessage(interact.ChannelID, messageID)
	if err != nil {
		return fmt.Errorf("GenerateByMessageComponent: %w", err)
	}

	if err := u.generate(ctx, interact.GuildID, interact.ChannelID, message.Author.ID, interact.Token, message); err != nil {
		return fmt.Errorf("GenerateByMessageComponent: %w", err)
	}

	return nil
}

func (u ImageUseCase) generate(ctx context.Context, guildID, channelID, userID, interactionToken string, message *discordgo.Message) error {
	var initImageURL *string
	var maskImageURL *string

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

			if len(message.Attachments) > 0 {
				maskImageURL = &message.Attachments[0].URL
			}
		}
	}

	if initImageURL == nil && len(message.Attachments) > 0 {
		initImageURL = &message.Attachments[0].URL
	}

	prompt, width, height, strength, err := resolveContent(message.Content)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	command := domain.ImageGenerateComamnd{
		Prompt:       prompt,
		Width:        width,
		Height:       height,
		Strength:     strength,
		InitImageURL: initImageURL,
		MaskImageURL: maskImageURL,
	}

	if err := u.imageService.Generate(ctx, command, imagenExtra(interactionToken, guildID, channelID, userID, message.ID)); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	return nil
}

func imagenExtra(interactionToken, guildID, channelID, userID, messageID string) map[string]interface{} {
	return map[string]interface{}{
		"via":               "discord",
		"user_id":           userID,
		"interaction_token": interactionToken,
		"message_id":        messageID,
		"message_url":       fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, messageID),
	}
}

func resolveContent(content string) (prompt string, width, height int, strength float64, err error) {
	var opt struct {
		Size     string  `short:"s"`
		Strength float64 `long:"strength"`
	}

	spl := strings.Split(content, "##")
	if len(spl) == 1 {
		return content, 0, 0, 0, nil
	}

	prompt = spl[0]
	optstr := spl[1]

	if err := commandline.ParseArgs(optstr, &opt); err != nil {
		return "", 0, 0, 0, fmt.Errorf("resolveContent: %w", err)
	}

	s := strings.Split(opt.Size, "x")

	width, err = strconv.Atoi(s[0])
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("resolveContent: %w", err)
	}

	height, err = strconv.Atoi(s[1])
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("resolveContent: %w", err)
	}

	strength = opt.Strength

	return
}
