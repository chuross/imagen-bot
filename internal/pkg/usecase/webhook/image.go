package webhook

import (
	"context"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/commandline"
	"imagen/internal/pkg/infra/discord"
	"imagen/internal/pkg/infra/environment"
	"imagen/internal/pkg/infra/service"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
)

type ImageUseCase struct {
	imageService domain.ImageService
}

func newImageUseCase(services *service.Services) *ImageUseCase {
	return &ImageUseCase{
		imageService: services.Image,
	}
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

	prompt, negativePrompts, width, height, strength, number, err := resolveContent(message.Content)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	command := domain.ImageGenerateComamnd{
		Prompt:          prompt,
		NegativePrompts: negativePrompts,
		RawPrompt:       message.Content,
		Width:           width,
		Height:          height,
		Strength:        strength,
		InitImageURL:    initImageURL,
		MaskImageURL:    maskImageURL,
		Number:          number,
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

func resolveContent(content string) (prompt string, negativePrompts []string, width, height int, strength float64, number int, err error) {
	var opt struct {
		Size     string  `long:"size"`
		Strength float64 `short:"s" long:"strength"`
		Number   int     `short:"n" default:"5"`
	}

	spl := strings.Split(content, "##")
	if len(spl) == 1 {
		nps := resolveNegativePrompts(spl[0])
		return content, nps, 0, 0, 0, 0, nil
	}

	prompt = spl[0]
	negativePrompts = resolveNegativePrompts(spl[0])
	optstr := spl[1]

	if err := commandline.ParseArgs(optstr, &opt); err != nil {
		return "", nil, 0, 0, 0, 0, fmt.Errorf("resolveContent: %w", err)
	}

	width, height, err = resolveSize(opt.Size)
	if err != nil {
		return "", nil, 0, 0, 0, 0, fmt.Errorf("resolveContent: %w", err)
	}

	strength = opt.Strength
	number = opt.Number

	return
}

func resolveSize(value string) (width, height int, err error) {
	s := strings.Split(value, "x")

	if len(s) != 2 {
		return 0, 0, nil
	}
	width, err = strconv.Atoi(s[0])
	if err != nil {
		return 0, 0, fmt.Errorf("resolveSize: %w", err)
	}

	height, err = strconv.Atoi(s[1])
	if err != nil {
		return 0, 0, fmt.Errorf("resolveSize: %w", err)
	}

	return
}

func resolveNegativePrompts(content string) []string {
	reg := regexp.MustCompile(`-\((.+)\)`).FindStringSubmatch(content)
	if len(reg) == 1 {
		return []string{}
	}

	ps := strings.Split(reg[1], ",")

	return lo.Map(ps, func(p string, _ int) string {
		return strings.TrimSpace(p)
	})
}
