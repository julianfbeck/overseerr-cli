package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	URL    string `json:"url"`
	APIKey string `json:"api_key"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "overseerr-cli", "config.json"), nil
}

func Load() (*Config, error) {
	// Check environment variables first
	url := os.Getenv("OVERSEERR_URL")
	apiKey := os.Getenv("OVERSEERR_API_KEY")

	if url != "" && apiKey != "" {
		return &Config{URL: url, APIKey: apiKey}, nil
	}

	// Fall back to config file
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Environment variables override config file
	if url != "" {
		cfg.URL = url
	}
	if apiKey != "" {
		cfg.APIKey = apiKey
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("OVERSEERR_URL not set. Use 'overseerr config set-url <url>' or set OVERSEERR_URL")
	}
	if c.APIKey == "" {
		return fmt.Errorf("OVERSEERR_API_KEY not set. Use 'overseerr config set-key <key>' or set OVERSEERR_API_KEY")
	}
	return nil
}

func SetURL(url string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.URL = url
	return cfg.Save()
}

func SetAPIKey(key string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.APIKey = key
	return cfg.Save()
}
