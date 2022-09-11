package middleware

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"imagen/internal/pkg/infra/environment"
	"log"
	"net/http"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func VerifyDiscordSignature(c *gin.Context) {
	log.Println("VerifyDiscordSignature: start")
	env := environment.MustGet(c.Request.Context())

	pubkey, err := hex.DecodeString(env.DISCORD.PUBLIC_KEY)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !discordgo.VerifyInteraction(c.Request, ed25519.PublicKey(pubkey)) {
		log.Println("VerifyDiscordSignature: end: 401")
		c.Status(http.StatusUnauthorized)
		return
	}
	log.Println("VerifyDiscordSignature: end: success")
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
