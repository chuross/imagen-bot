package middleware

import (
	"crypto/ed25519"
	"imagen/internal/pkg/infra/environment"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func VerifyDiscordSignature(c *gin.Context) {
	env := environment.MustGet(c.Request.Context())

	if !discordgo.VerifyInteraction(c.Request, ed25519.PublicKey(env.DISCORD.PUBLIC_KEY)) {
		c.Status(http.StatusUnauthorized)
		return
	}
}
