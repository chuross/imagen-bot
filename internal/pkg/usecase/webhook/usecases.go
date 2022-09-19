package webhook

import "imagen/internal/pkg/infra/service"

type UseCases struct {
	Image     *ImageUseCase
	Workspace *WorkspaceUseCase
}

func NewWebhookUseCases(services *service.Services) *UseCases {
	return &UseCases{
		Image:     newImageUseCase(services),
		Workspace: newWorkspaceUseCase(services),
	}
}
