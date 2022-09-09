package main

import (
	"imagen/internal/pkg/interface/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.Setup(r)

	r.Run()
}
