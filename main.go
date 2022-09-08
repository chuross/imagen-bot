package main

import (
	"fmt"
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
	if ctx, err := env.With(c.Request.Context()); err != nil {
		fmt.Printf("Error: initialize env: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.Request = c.Request.WithContext(ctx)
	}
}

func webhook(c *gin.Context) {
	c.Status(http.StatusOK)
}
