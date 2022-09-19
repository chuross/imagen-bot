package webhook

import (
	"fmt"
	"imagen/internal/pkg/domain"
	"imagen/internal/pkg/infra/service"
)

type WorkspaceUseCase struct {
	workspaceService domain.WorkspaceService
}

func newWorkspaceUseCase(services *service.Services) *WorkspaceUseCase {
	return &WorkspaceUseCase{
		workspaceService: services.Workspace,
	}
}

func (w WorkspaceUseCase) Create(channelID, messageID string) error {
	if err := w.workspaceService.Create(channelID, messageID); err != nil {
		return fmt.Errorf("Create: %w", err)
	}
	return nil
}
