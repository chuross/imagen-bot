package main

import (
	"fmt"
	"net/http"

	"imagen/internal/pkg/infra/environment"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	r := gin.Default()

	root := r.Use(withEnv)
	root.POST("/hook", webhook)

	r.Run()
}

func withEnv(c *gin.Context) {
	if ctx, err := environment.With(c.Request.Context()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.Request = c.Request.WithContext(ctx)
	}
}

func webhook(c *gin.Context) {
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
