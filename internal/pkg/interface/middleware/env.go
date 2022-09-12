package middleware

import (
	"imagen/internal/pkg/infra/environment"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WithEnv(c *gin.Context) {
	if ctx, err := environment.With(c.Request.Context()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.Request = c.Request.WithContext(ctx)
	}
}
