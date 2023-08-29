package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/urfave/cli/v2"
)

func handleCLICurl(ctx *cli.Context) error {
	profile := getCLIProfile(ctx.Context)

	// We have to extract the curl arguments from os.Args because urfave has intercepted
	// and removed the flags.
	args := os.Args
	if curlIdx := slices.Index(args, "curl"); curlIdx != -1 && curlIdx < len(args)-1 {
		args = args[curlIdx+1:]
	} else {
		return fmt.Errorf("usage: nola curl [curl arguments]")
	}

	cmd, err := executeCurl(profile.AccessToken, args)
	if err != nil {
		return err
	}

	os.Exit(cmd.ProcessState.ExitCode())

	return nil
}
