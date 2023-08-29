package main

import (
	"context"

	"nolahq.com/cli/config"
)

type contextKey int

const (
	contextKeyConfig contextKey = iota
	contextKeyProfile
)

func createCLIContext(ctx context.Context, conf *config.Config) context.Context {
	return context.WithValue(ctx, contextKeyConfig, conf)
}

func getCLIConfig(ctx context.Context) *config.Config {
	v, ok := ctx.Value(contextKeyConfig).(*config.Config)
	if !ok {
		panic("config missing in context, this should not happen")
	}
	return v
}

func createCLIProfileContext(ctx context.Context, profile *config.Profile) context.Context {
	return context.WithValue(ctx, contextKeyProfile, profile)
}

func getCLIProfile(ctx context.Context) *config.Profile {
	v, ok := ctx.Value(contextKeyProfile).(*config.Profile)
	if !ok {
		panic("profile missing in context, this should not happen")
	}
	return v
}
