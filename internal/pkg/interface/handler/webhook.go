package handler

import (
	"fmt"
	"imagen/internal/pkg/infra/discord"
	"imagen/internal/pkg/usecase/webhook"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	imageUseCase     *webhook.ImageUseCase
	workspaceUseCase *webhook.WorkspaceUseCase
}

func NewWebhookHandler(usecases *webhook.UseCases) *WebhookHandler {
	return &WebhookHandler{
		imageUseCase:     usecases.Image,
		workspaceUseCase: usecases.Workspace,
	}
}

func (h WebhookHandler) HookByDiscord(c *gin.Context) {
	intaract := discordgo.Interaction{}

	if err := c.BindJSON(&intaract); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("HookByDiscord: %w", err))
		return
	}

	switch intaract.Type {
	case discordgo.InteractionApplicationCommand:
		data := intaract.ApplicationCommandData()

		log.Printf("HookByDiscord: receive application command: name=%v", data.Name)

		switch data.Name {
		case discord.CommandImagen.Name:
			if err := h.imageUseCase.GenerateByMessageCommand(c.Request.Context(), &intaract); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		case discord.CommandWorkspace.Name:
			if err := h.workspaceUseCase.Create()
		}

		c.JSON(http.StatusOK, gin.H{
			"type": discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	case discordgo.InteractionMessageComponent:
		log.Println("HookByDiscord: receive message component command")

		if err := h.imageUseCase.GenerateByMessageComponent(c.Request.Context(), &intaract); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"type": discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	default:
		c.JSON(http.StatusOK, gin.H{
			"type": discordgo.InteractionResponsePong,
		})
	}
}
