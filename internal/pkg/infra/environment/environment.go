package environment

import (
	"log"

	"github.com/Netflix/go-env"
)

var environ = &Env{}

type Env struct {
	GOOGLE_CLOUD_PROJECT_ID string `env:"GOOGLE_CLOUD_PROJECT_ID,required=true"`

	LINE_BOT struct {
		CHANNEL_ACCESS_TOKEN string `env:"LINE_BOT_CHANNEL_ACCESS_TOKEN,required=true"`
		SECRET_TOKEN         string `env:"LINE_BOT_SECRET_TOKEN,reuqired=true"`
	}

	DISCORD struct {
		PUBLIC_KEY string `env:"DISCORD_PUBLIC_KEY,required=true"`
		BOT_TOKEN  string `env:"DISCORD_BOT_TOKEN,required=true"`
		APP_ID     string `env:"DISCORD_APP_ID,required=true"`
		GUILD_ID   string `env:"DISCORD_GUILD_ID,required=true"`
	}
}

func MustGet() *Env {
	if environ == nil {
		panic("env is not set! must call Load")
	}
	return environ
}

func Load() error {
	log.Printf("load environment")
	if _, err := env.UnmarshalFromEnviron(environ); err != nil {
		return err
	}
	return nil
}
