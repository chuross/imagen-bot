package service

import (
	"imagen/internal/pkg/domain"
)

type Services struct {
	Image     domain.ImageService
	Workspace domain.WorkspaceService
}

func NewServices() *Services {
	return &Services{
		Image:     newImageService(),
		Workspace: newWorkspaceService(),
	}
}
