package router

import (
	"imagen/internal/pkg/interface/handler"
	"imagen/internal/pkg/interface/middleware"
	"imagen/internal/pkg/usecase/webhook"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, webhookUseCases *webhook.UseCases) {
	root := r.Group("/")
	root.Use(gin.Recovery(), middleware.WithEnv)

	{
		handler := handler.NewWebhookHandler(webhookUseCases)
		root.POST("/hooks/line", handler.HookByLine)

		discord := root.Group("")
		discord.Use(middleware.VerifyDiscordSignature)
		{
			discord.POST("/hooks/discord", handler.HookByDiscord)
		}
	}
}
