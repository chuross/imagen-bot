package handler

import (
	"fmt"
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

func (h WebhookHandler) Hook(c *gin.Context) {
	e := environment.MustGet(c.Request.Context())

	bot, err := linebot.New(e.LINE_BOT.SECRET_TOKEN, e.LINE_BOT.CHANNEL_ACCESS_TOKEN)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	for _, event := range events {
		fmt.Printf("event received: type=%v, messageType=%v", event.Type, event.Message.Type())

		if event.Type != linebot.EventTypeMessage {
			continue
		}

		if event.Message.Type() != linebot.MessageTypeText {
			continue
		}

		if err := h.imageUseCase.Generate(c.Request.Context(), event.Message.(*linebot.TextMessage).Text); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		mes := linebot.NewTextMessage(messageSuccess)
		if _, err := bot.ReplyMessage(event.ReplyToken, mes).Do(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	c.Status(http.StatusOK)
}
