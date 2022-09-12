package handler

import (
	"fmt"
	"imagen/internal/pkg/infra/environment"
	"imagen/internal/pkg/usecase/webhook"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type WebhookHandler struct {
	imageUseCase *webhook.ImageUseCase
}

func NewWebhookHandler(usecases *webhook.UseCases) *WebhookHandler {
	return &WebhookHandler{
		imageUseCase: usecases.Image,
	}
}

func (h WebhookHandler) HookByLine(c *gin.Context) {
	e := environment.MustGet(c.Request.Context())

	bot, err := linebot.New(e.LINE_BOT.SECRET_TOKEN, e.LINE_BOT.CHANNEL_ACCESS_TOKEN)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := h.imageUseCase.GenerateByLine(c.Request.Context(), events); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
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
		if err := h.imageUseCase.GenerateByDiscord(c.Request.Context(), &intaract); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"type": discordgo.InteractionResponseChannelMessageWithSource,
			"data": gin.H{
				"content": "生成リクエストを経由します...",
			},
		})
	default:
		c.JSON(http.StatusOK, gin.H{
			"type": discordgo.InteractionResponsePong,
		})
	}
}
