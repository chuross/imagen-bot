package webhook

import (
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/discord"
	"imagen/internal/pkg/infra/service"

	"github.com/bwmarrin/discordgo"
)

type WorkspaceUseCase struct {
	workspaceService domain.WorkspaceService
}

func newWorkspaceUseCase(services *service.Services) *WorkspaceUseCase {
	return &WorkspaceUseCase{
		workspaceService: services.Workspace,
	}
}

func (w WorkspaceUseCase) Create(interact *discordgo.Interaction) error {
	data := interact.ApplicationCommandData()

	if data.Name != discord.CommandWorkspace.Name {
		return fmt.Errorf("Create: unexpected command: %v", interact.ApplicationCommandData().Name)
	}

	message, ok := data.Resolved.Messages[data.TargetID]
	if !ok {
		return nil
	}

	if err := w.workspaceService.Create(interact.ChannelID, message.ID); err != nil {
		return fmt.Errorf("Create: %w", err)
	}

	return nil
}
