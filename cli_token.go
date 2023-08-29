package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"nolahq.com/cli/auth"
	"nolahq.com/cli/config"
)

func handleCLIGetToken(ctx *cli.Context) error {
	profile := getCLIProfile(ctx.Context)

	fmt.Println(profile.AccessToken)
	os.Exit(0)
	return nil
}

func refreshCredentialsIfNecessary(ctx *cli.Context) error {
	profileName, profile, err := findProfile(ctx)
	if err != nil {
		return err
	}

	if profile.Expiry.After(time.Now()) {
		return nil
	}

	// If we're here, the access token has expired and we need to
	// refresh it.
	creds := &auth.Credentials{
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
	}

	if err := creds.Refresh(profile.Server); err != nil {
		return fmt.Errorf("failed to refresh token: %w. try `nola login` from scratch.", err)
	}

	profile.AccessToken = creds.AccessToken
	profile.RefreshToken = creds.RefreshToken
	profile.Expiry = creds.Expiry

	conf := getCLIConfig(ctx.Context)
	conf.AddProfile(
		profileName, profile.Server, profile.Principal,
		profile.AccessToken, profile.RefreshToken, profile.Expiry,
	)
	if err := conf.Save(); err != nil {
		return fmt.Errorf("failed to save the refreshed credentials for profile %s: %w", profileName, err)
	}

	return nil
}

func findProfile(ctx *cli.Context) (string, *config.Profile, error) {
	conf := getCLIConfig(ctx.Context)

	profileName := ctx.String("profile")
	if profileName == "" {
		profileName = fmt.Sprintf("default-%s", ctx.String("server"))
	}

	profile := conf.GetProfile(profileName)
	if profile == nil {
		return "", nil, fmt.Errorf("profile `%s` not found. Try `nola login --help`.", profileName)
	}

	return profileName, profile, nil

}
