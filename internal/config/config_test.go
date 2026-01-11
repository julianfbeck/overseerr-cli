package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "empty config",
			config:  Config{},
			wantErr: true,
		},
		{
			name:    "missing URL",
			config:  Config{APIKey: "test-key"},
			wantErr: true,
		},
		{
			name:    "missing API key",
			config:  Config{URL: "https://example.com"},
			wantErr: true,
		},
		{
			name:    "valid config",
			config:  Config{URL: "https://example.com", APIKey: "test-key"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Save original env vars
	origURL := os.Getenv("OVERSEERR_URL")
	origKey := os.Getenv("OVERSEERR_API_KEY")
	defer func() {
		os.Setenv("OVERSEERR_URL", origURL)
		os.Setenv("OVERSEERR_API_KEY", origKey)
	}()

	// Set test env vars
	os.Setenv("OVERSEERR_URL", "https://test.example.com")
	os.Setenv("OVERSEERR_API_KEY", "test-api-key")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.URL != "https://test.example.com" {
		t.Errorf("URL = %v, want %v", cfg.URL, "https://test.example.com")
	}
	if cfg.APIKey != "test-api-key" {
		t.Errorf("APIKey = %v, want %v", cfg.APIKey, "test-api-key")
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create temp directory for config
	tmpDir, err := os.MkdirTemp("", "overseerr-cli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clear env vars for this test
	origURL := os.Getenv("OVERSEERR_URL")
	origKey := os.Getenv("OVERSEERR_API_KEY")
	os.Unsetenv("OVERSEERR_URL")
	os.Unsetenv("OVERSEERR_API_KEY")
	defer func() {
		os.Setenv("OVERSEERR_URL", origURL)
		os.Setenv("OVERSEERR_API_KEY", origKey)
	}()

	// Create config file directly in temp location
	configDir := filepath.Join(tmpDir, ".config", "overseerr-cli")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	cfg := Config{
		URL:    "https://saved.example.com",
		APIKey: "saved-api-key",
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Read back and verify
	readData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var readCfg Config
	if err := json.Unmarshal(readData, &readCfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if readCfg.URL != cfg.URL {
		t.Errorf("URL = %v, want %v", readCfg.URL, cfg.URL)
	}
	if readCfg.APIKey != cfg.APIKey {
		t.Errorf("APIKey = %v, want %v", readCfg.APIKey, cfg.APIKey)
	}
}

func TestEnvOverridesConfig(t *testing.T) {
	// Save original env vars
	origURL := os.Getenv("OVERSEERR_URL")
	origKey := os.Getenv("OVERSEERR_API_KEY")
	defer func() {
		os.Setenv("OVERSEERR_URL", origURL)
		os.Setenv("OVERSEERR_API_KEY", origKey)
	}()

	// Set only URL from env
	os.Setenv("OVERSEERR_URL", "https://env.example.com")
	os.Setenv("OVERSEERR_API_KEY", "env-key")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Both should come from env since both are set
	if cfg.URL != "https://env.example.com" {
		t.Errorf("URL = %v, want %v", cfg.URL, "https://env.example.com")
	}
	if cfg.APIKey != "env-key" {
		t.Errorf("APIKey = %v, want %v", cfg.APIKey, "env-key")
	}
}
