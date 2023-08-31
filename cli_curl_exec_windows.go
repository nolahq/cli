//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const curlTpl = `
@echo off
SET NOLA_SECRET={{ .Token }}
curl -H "Authorization: Bearer %NOLA_SECRET%" {{ .Args }}
`

func executeCurl(bearerToken string, curlArgs []string) (*exec.Cmd, error) {
	f, err := os.CreateTemp("", "nola-cli-curl-*.bat")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file for curl: %w", err)
	}

	defer f.Close()
	defer os.Remove(f.Name())

	if err := os.Chmod(f.Name(), 0700); err != nil {
		return nil, fmt.Errorf("failed to chmod temp file for curl: %w", err)
	}

	tpl := template.Must(template.New("").Parse(curlTpl))
	if err := tpl.Execute(f, map[string]any{
		"Token": bearerToken,
		"Args":  strings.Join(curlArgs, " "),
	}); err != nil {
		return nil, fmt.Errorf("failed to execute curl template: %w", err)
	}
	if err := f.Sync(); err != nil {
		return nil, fmt.Errorf("failed to sync temp file for curl: %w", err)
	}

	cmd := exec.Command("cmd.exe", "/C", f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute curl: %w", err)
	}

	return cmd, nil
}
