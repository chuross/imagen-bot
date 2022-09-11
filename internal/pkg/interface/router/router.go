package router

import (
	"imagen/internal/pkg/interface/handler"
	"imagen/internal/pkg/interface/middleware"
	"imagen/internal/pkg/usecase/webhook"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, webhookUseCases *webhook.UseCases) {
	root := r.Use(gin.Recovery()).
		Use(middleware.WithEnv)

	{
		handler := handler.NewWebhookHandler(webhookUseCases)
		root.POST("/hooks/line", handler.HookByLine)
	}
}
