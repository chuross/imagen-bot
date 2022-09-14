package router

import (
	"imagen/internal/pkg/interface/handler"
	"imagen/internal/pkg/interface/middleware"
	"imagen/internal/pkg/usecase/webhook"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, webhookUseCases *webhook.UseCases) {
	root := r.Group("/", gin.Recovery())

	{
		handler := handler.NewWebhookHandler(webhookUseCases)
		root.POST("/hooks/discord", middleware.VerifyDiscordSignature, middleware.RegisterInteractionCommand, handler.HookByDiscord)
	}
}
