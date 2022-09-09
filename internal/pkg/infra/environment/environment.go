package environment

import (
	"context"
	"fmt"

	"github.com/Netflix/go-env"
)

type envKey struct{}

var envKeyMain = envKey{}

type Env struct {
	GOOGLE_CLOUD_PROJECT_ID string

	LINE_BOT struct {
		CHANNEL_ACCESS_TOKEN string `env:"CHANNEL_ACCESS_TOKEN,required=true"`
		SECRET_TOKEN         string `env:"SECRET_TOKEN,reuqired=true"`
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
		return ctx, fmt.Errorf("unmarshal environment failed: %w", err)
	} else {
		return context.WithValue(ctx, envKeyMain, &e), nil
	}
}
