package router

import (
	"imagen/internal/pkg/interface/handler"
	"imagen/internal/pkg/interface/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	root := r.Use(middleware.WithEnv)
	root.POST("/hook", handler.Webhook)
}
