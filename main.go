package main

import (
	"imagen/internal/pkg/infra/imagen/service"
	"imagen/internal/pkg/interface/router"
	"imagen/internal/pkg/usecase/webhook"

	"github.com/gin-gonic/gin"
)

func main() {
	services := service.NewServices()

	r := gin.Default()

	router.Setup(
		r,
		webhook.NewWebhookUseCases(services),
	)

	r.Run()
}
