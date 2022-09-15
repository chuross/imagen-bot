package handler

import (
	"fmt"
	"imagen/internal/pkg/usecase/webhook"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	imageUseCase *webhook.ImageUseCase
}

func NewWebhookHandler(usecases *webhook.UseCases) *WebhookHandler {
	return &WebhookHandler{
		imageUseCase: usecases.Image,
	}
}

func (h WebhookHandler) HookByDiscord(c *gin.Context) {
	intaract := discordgo.Interaction{}

	if err := c.BindJSON(&intaract); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("HookByDiscord: %w", err))
		return
	}

	log.Printf("HookByDiscord: receive interaction: type=%v", intaract.Type)

	switch intaract.Type {
	case discordgo.InteractionApplicationCommand:
		if err := h.imageUseCase.GenerateByMessageCommand(c.Request.Context(), &intaract); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"type": discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	case discordgo.InteractionMessageComponent:
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
