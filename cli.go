package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/urfave/cli/v2"
	"nolahq.com/cli/config"
)

func main() {
	prepareCredentialsAction := func(ctx *cli.Context) error {
		if err := refreshCredentialsIfNecessary(ctx); err != nil {
			return fmt.Errorf("credentials pre-check failed: %w", err)
		}

		_, profile, err := findProfile(ctx)
		if err != nil {
			return fmt.Errorf("failed to find profile: %w", err)
		}

		ctx.Context = createCLIProfileContext(ctx.Context, profile)
		return nil
	}

	app := &cli.App{
		Name:  "nola",
		Usage: "A command line tool for interacting with https://app.nolahq.com",
		Before: func(ctx *cli.Context) error {
			conf, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			ctx.Context = createCLIContext(ctx.Context, conf)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "profile",
				Aliases:     []string{"p"},
				Usage:       "The name of the profile to load/save",
				DefaultText: "",
			},
			&cli.StringFlag{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "The Nola server to connect to (live, staging, dev)",
				Value:   "live",
				Action: func(ctx *cli.Context, s string) error {
					if !slices.Contains([]string{"live", "staging", "dev"}, s) {
						return errors.New("valid choices for --server are: live, staging, dev")
					}
					ctx.Set("server", s)
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "login",
				Usage:  "Login to Nola and save the credentials to the config file",
				Action: handleCLILogin,
			},
			{
				Name:   "token",
				Usage:  "Get an access token suitable for REST authentication",
				Action: handleCLIGetToken,
				Before: prepareCredentialsAction,
			},
			{
				Name:            "curl",
				Usage:           "Run a curl command with authentication pre-filled.",
				Action:          handleCLICurl,
				Before:          prepareCredentialsAction,
				SkipFlagParsing: true,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getProfileName(ctx *cli.Context) string {
	profileName := ctx.String("profile")
	if profileName == "" {
		profileName = fmt.Sprintf("default-%s", ctx.String("server"))
	}
	return profileName
}
