package env

import (
	"context"
	"fmt"

	"github.com/Netflix/go-env"
)

type envKey struct{}

var envKeyMain = envKey{}

type Env struct {
	LINE_BOT struct {
		CHANNEL_ACCESS_TOKEN string `env:"CHANNEL_ACCESS_TOKEN,required=true"`
		SECRET_TOKEN         string `env:"SECRET_TOKEN,reuqired=true"`
	}
}

func Get(ctx context.Context) Env {
	e := ctx.Value(envKeyMain).(*Env)
	if e == nil {
		panic("env is not set! must call WithEnv")
	}
	return *e
}

func WithEnv(ctx context.Context) (context.Context, error) {
	var environment Env
	if _, err := env.UnmarshalFromEnviron(&environment); err != nil {
		return ctx, fmt.Errorf("unmarshal environment failed: %w", err)
	} else {
		return context.WithValue(ctx, envKeyMain, &environment), nil
	}
}
