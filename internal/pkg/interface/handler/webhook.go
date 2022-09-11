package handler

import (
	"imagen/internal/pkg/infra/environment"
	"imagen/internal/pkg/usecase/webhook"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	messageSuccess = "画像リクエスト受け付けました"
)

type WebhookHandler struct {
	imageUseCase webhook.ImageUseCase
}

func NewWebhookHandler(usecases webhook.UseCases) WebhookHandler {
	return WebhookHandler{
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
