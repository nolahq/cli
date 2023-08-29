package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cli/oauth"
)

type Credentials struct {
	AccessToken  string
	RefreshToken string
	Principal    string
	Expiry       time.Time
}

func Authenticate(server string) (*Credentials, error) {
	flow := &oauth.Flow{
		Host:     oauthHost(server),
		ClientID: "nola-cli",
		Scopes:   []string{"profile", "email", "offline_access", "openid"},
	}

	result, err := flow.DetectFlow()
	if err != nil {
		return nil, fmt.Errorf("authorization failed: %w", err)
	}

	if strings.ToLower(result.Type) != "bearer" {
		return nil, fmt.Errorf("expected bearer token, got %s", result.Type)
	}

	email, err := tokenToPrincipal(result.Token, server)
	if err != nil {
		return nil, fmt.Errorf("couldn't determine principal email: %w", err)
	}

	creds := &Credentials{
		AccessToken:  result.Token,
		RefreshToken: result.RefreshToken,
		Principal:    email,
		Expiry:       time.Now(),
	}

	if err := creds.Refresh(server); err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w. try `nola login` from scratch.", err)
	}

	return creds, nil
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

func (c *Credentials) Refresh(server string) error {
	u := fmt.Sprintf(
		"https://auth.nolahq.com/realms/%s/protocol/openid-connect/token",
		serverToRealm(server),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", c.RefreshToken)
	v.Set("client_id", "nola-cli")

	req, _ := http.NewRequest(http.MethodPost, u, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)

	now := time.Now()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("token refresh failed to send: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token refresh failed with status %d", resp.StatusCode)
	}

	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("error decoding auth refresh_token response: %w", err)
	}

	if strings.ToLower(body.TokenType) != "bearer" {
		return fmt.Errorf("expected bearer token, got %s", body.TokenType)
	}

	c.AccessToken = body.AccessToken
	c.Expiry = now.Add(time.Duration(body.ExpiresIn) * time.Second)

	return nil
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

func oauthHost(server string) *oauth.Host {
	realm := serverToRealm(server)
	return &oauth.Host{
		DeviceCodeURL: fmt.Sprintf(
			"https://auth.nolahq.com/realms/%s/protocol/openid-connect/auth/device", realm),
		AuthorizeURL: fmt.Sprintf(
			"https://auth.nolahq.com/realms/%s/protocol/openid-connect/auth", realm),
		TokenURL: fmt.Sprintf(
			"https://auth.nolahq.com/realms/%s/protocol/openid-connect/token", realm),
	}
}
