package handler

import (
	"fmt"
	"imagen/internal/pkg/infra/environment"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func Webhook(c *gin.Context) {
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
		if event.Type != linebot.EventTypeMessage {
			continue
		}
		if event.Message.Type() != linebot.MessageTypeText {
			continue
		}
		fmt.Println(event)
	}

	c.Status(http.StatusOK)
}
