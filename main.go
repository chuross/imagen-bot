package main

import (
	"imagen/pkg/infrastructure/env"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	root := r.Use(withEnv)
	root.POST("/hook", webhook)

	r.Run()
}

func withEnv(c *gin.Context) {
	env.WithEnv(c.Request.Context())
}

func webhook(c *gin.Context) {
	c.Status(http.StatusOK)
}
