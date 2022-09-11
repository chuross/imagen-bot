package middleware

import (
	"crypto/ed25519"
	"fmt"
	"imagen/internal/pkg/infra/environment"
	"net/http"
	"sync"

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

func RegisterInteractionCommand(c *gin.Context) {
	env := environment.MustGet(c.Request.Context())
	if _, err := discordgo.New(fmt.Sprintf("Bot %s", env.DISCORD.BOT_TOKEN)); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	} else {
		var once sync.Once
		once.Do(func() {
		})
	}

}