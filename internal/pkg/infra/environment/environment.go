package environment

import (
	"context"

	"github.com/Netflix/go-env"
)

type envKey struct{}

var envKeyMain = envKey{}

type Env struct {
	GOOGLE_CLOUD_PROJECT_ID string `env:"GOOGLE_CLOUD_PROJECT_ID,required=true"`

	LINE_BOT struct {
		CHANNEL_ACCESS_TOKEN string `env:"LINE_BOT_CHANNEL_ACCESS_TOKEN,required=true"`
		SECRET_TOKEN         string `env:"LINE_BOT_SECRET_TOKEN,reuqired=true"`
	}

	DISCORD struct {
		PUBLIC_KEY string `env:"DISCORD_PUBLIC_KEY,required=true"`
		BOT_TOKEN  string `env:"DISCORD_BOT_TOKEN,required=true"`
	}
}

func MustGet(ctx context.Context) Env {
	env := ctx.Value(envKeyMain).(*Env)
	if env == nil {
		panic("env is not set! must call WithEnv")
	}
	return *env
}

func With(ctx context.Context) (context.Context, error) {
	var e Env
	if _, err := env.UnmarshalFromEnviron(&e); err != nil {
		return ctx, err
	} else {
		return context.WithValue(ctx, envKeyMain, &e), nil
	}
}
