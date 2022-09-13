package main

import (
	"imagen/internal/pkg/infra/environment"
	"imagen/internal/pkg/infra/service"
	"imagen/internal/pkg/interface/router"
	"imagen/internal/pkg/usecase/webhook"

	"github.com/gin-gonic/gin"
)

func main() {
	environment.Load()
	services := service.NewServices()

	r := gin.Default()

	router.Setup(
		r,
		webhook.NewWebhookUseCases(services),
	)

	r.Run()
}
