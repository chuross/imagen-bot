package service

import (
	"bytes"
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/environment"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
)

type workspaceService struct{}

func newWorkspaceService() domain.WorkspaceService {
	return &workspaceService{}
}

func (w workspaceService) Create(channelID, messageID string) error {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", environment.MustGet().DISCORD.BOT_TOKEN))
	if err != nil {
		return fmt.Errorf("Create: %w", err)
	}

	message, err := session.ChannelMessage(channelID, messageID)
	if err != nil {
		return fmt.Errorf("Create: %w", err)
	}

	now := time.Now()
	threadName := fmt.Sprintf("workspace-%s", now.Format("20060102"))

	ch, err := session.ThreadStart(channelID, threadName, discordgo.ChannelTypeGuildPublicThread, 60*24*3)
	if err != nil {
		return fmt.Errorf("Create: %w", err)
	}

	attachment := message.Attachments[0]

	r, err := resty.New().R().Get(attachment.URL)
	if err != nil {
		return fmt.Errorf("Create: %w", err)
	}

	if _, err := session.ChannelFileSendWithMessage(ch.ID, message.Content, attachment.Filename, bytes.NewReader(r.Body())); err != nil {
		return fmt.Errorf("Create: %w", err)
	}

	return nil
}
