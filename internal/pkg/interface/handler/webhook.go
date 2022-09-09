package handler

import (
	"fmt"
	"imagen/internal/pkg/infra/environment"
	"imagen/internal/pkg/usecase/webhook"
	"log"
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
		return
	}

	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, event := range events {
		if event.Type != linebot.EventTypeMessage {
			log.Printf("unexpected event type: type=%v", event.Type)
			continue
		}

		switch event.Message.(type) {
		case *linebot.TextMessage:
			text := event.Message.(*linebot.TextMessage).Text
			sendingTargetID := event.Source.UserID

			log.Printf("generate image: text=%v, sendingTargetID=%v", text, sendingTargetID)
			if err := h.imageUseCase.Generate(c.Request.Context(), text, sendingTargetID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			mes := linebot.NewTextMessage(messageSuccess)
			if _, err := bot.ReplyMessage(event.ReplyToken, mes).Do(); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		default:
			fmt.Printf("unexpected message type: type=%v", event.Message)
		}
	}

	c.Status(http.StatusOK)
}
