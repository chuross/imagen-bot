package domain

type WorkspaceService interface {
	Create(channelID, messageID string) error
}
