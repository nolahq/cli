package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cli/oauth"
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

func tokenToPrincipal(token, server string) (string, error) {
	u := fmt.Sprintf(
		"https://auth.nolahq.com/realms/%s/protocol/openid-connect/userinfo",
		serverToRealm(server),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var body struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("error decoding auth user_endpoint response: %w", err)
	}

	if body.Email == "" {
		return "", fmt.Errorf("no email in auth user_endpoint response")
	}

	return body.Email, nil
}

func serverToRealm(server string) string {
	realm := "nola"
	switch server {
	case "staging":
		realm = "nola-staging"
	case "dev":
		realm = "nola-localhost"
	}
	return realm
}

func oauthHost(ctx *cli.Context) *oauth.Host {
	realm := serverToRealm(ctx.String("server"))
	return &oauth.Host{
		DeviceCodeURL: fmt.Sprintf(
			"https://auth.nolahq.com/realms/%s/protocol/openid-connect/auth/device", realm),
		AuthorizeURL: fmt.Sprintf(
			"https://auth.nolahq.com/realms/%s/protocol/openid-connect/auth", realm),
		TokenURL: fmt.Sprintf(
			"https://auth.nolahq.com/realms/%s/protocol/openid-connect/token", realm),
	}
}
