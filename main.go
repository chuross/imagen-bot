package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/hook", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	r.Run()
}
