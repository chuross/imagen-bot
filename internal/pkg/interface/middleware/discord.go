package middleware

import (
	"crypto/ed25519"
	"encoding/hex"
	"imagen/internal/pkg/infra/discord"
	"imagen/internal/pkg/infra/environment"
	"log"
	"net/http"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func VerifyDiscordSignature(c *gin.Context) {
	env := environment.MustGet()

	pubkey, err := hex.DecodeString(env.DISCORD.PUBLIC_KEY)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !discordgo.VerifyInteraction(c.Request, ed25519.PublicKey(pubkey)) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func RegisterInteractionCommand(c *gin.Context) {
	env := environment.MustGet()

	var once sync.Once
	once.Do(func() {
		ses, err := discordgo.New("Bot " + env.DISCORD.BOT_TOKEN)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		appID := env.DISCORD.APP_ID
		guildID := env.DISCORD.GUILD_ID

		if _, err := ses.ApplicationCommandBulkOverwrite(appID, guildID, discord.Commands); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		log.Println("RegisterInteractionCommand: all commands updated")
	})
}
