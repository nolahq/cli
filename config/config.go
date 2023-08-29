package config

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Profiles map[string]Profile `toml:"profiles"`
}

type Profile struct {
	Principal    string    `toml:"principal"`
	AccessToken  string    `toml:"access_token"`
	RefreshToken string    `toml:"refresh_token"`
	Server       string    `toml:"server"`
	Expiry       time.Time `toml:"expiry"`
}

func (c *Config) AddProfile(
	profileName, server,
	principal, accessToken, refreshToken string, expiry time.Time) *Profile {
	profile := Profile{
		Principal:    principal,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Server:       server,
		Expiry:       expiry,
	}
	if c.Profiles == nil {
		c.Profiles = make(map[string]Profile)
	}
	c.Profiles[profileName] = profile

	return &profile
}

func (c *Config) GetProfile(profileName string) *Profile {
	if c.Profiles == nil {
		return nil
	}
	profile, ok := c.Profiles[profileName]
	if !ok {
		return nil
	}
	return &profile
}

func Load() (*Config, error) {
	configPath := getConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get config file stat at %s: %w", configPath, err)
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file %s: %w", configPath, err)
	}

	return &config, nil
}

func (c *Config) Save() error {
	configPath := getConfigPath()

	if err := os.MkdirAll(path.Dir(configPath), 0700); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", path.Dir(configPath), err)
	}

	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", configPath, err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(c); err != nil {
		return fmt.Errorf("failed to encode config file %s: %w", configPath, err)
	}

	return nil
}

func getConfigPath() string {
	dir, _ := os.UserConfigDir()
	return path.Join(dir, "com.nolahq", "cli.toml")
}
