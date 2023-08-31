package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"nolahq.com/cli/auth"
)

func handleCLILogin(ctx *cli.Context) error {
	server := ctx.String("server")

	creds, err := auth.Authenticate(server)
	if err != nil {
		return err
	}

	conf := getCLIConfig(ctx.Context)
	conf.AddProfile(
		getProfileName(ctx), server,
		creds.Principal, creds.AccessToken, creds.RefreshToken, creds.Expiry,
	)
	if err := conf.Save(); err != nil {
		return fmt.Errorf("failed to save the new profile: %w", err)
	}

	return nil
}
