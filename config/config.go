// Package config provides configuration management for sub2api.
// It handles loading and validating application settings from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration.
type Config struct {
	// Server settings
	Host string
	Port int

	// Subscription settings
	SubURL      string
	RefreshInterval int // in minutes

	// API settings
	APIToken    string
	BasePath    string

	// Output settings
	OutputFormat string // clash, singbox, raw
	NodeFilter   string

	// Logging
	LogLevel string
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Host:            "0.0.0.0",
		Port:            8080,
		RefreshInterval: 30,
		BasePath:        "/",
		OutputFormat:    "clash",
		LogLevel:        "info",
	}
}

// LoadFromEnv loads configuration from environment variables,
// falling back to defaults where values are not set.
func LoadFromEnv() (*Config, error) {
	cfg := DefaultConfig()

	if v := os.Getenv("HOST"); v != "" {
		cfg.Host = v
	}

	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT value %q: %w", v, err)
		}
		cfg.Port = p
	}

	if v := os.Getenv("SUB_URL"); v != "" {
		cfg.SubURL = v
	}

	if v := os.Getenv("REFRESH_INTERVAL"); v != "" {
		ri, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid REFRESH_INTERVAL value %q: %w", v, err)
		}
		cfg.RefreshInterval = ri
	}

	if v := os.Getenv("API_TOKEN"); v != "" {
		cfg.APIToken = v
	}

	if v := os.Getenv("BASE_PATH"); v != "" {
		// Ensure base path starts with /
		if !strings.HasPrefix(v, "/") {
			v = "/" + v
		}
		cfg.BasePath = v
	}

	if v := os.Getenv("OUTPUT_FORMAT"); v != "" {
		cfg.OutputFormat = strings.ToLower(v)
	}

	if v := os.Getenv("NODE_FILTER"); v != "" {
		cfg.NodeFilter = v
	}

	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = strings.ToLower(v)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required fields are set and values are within acceptable ranges.
func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("PORT must be between 1 and 65535, got %d", c.Port)
	}

	if c.RefreshInterval < 1 {
		return fmt.Errorf("REFRESH_INTERVAL must be at least 1 minute, got %d", c.RefreshInterval)
	}

	validFormats := map[string]bool{"clash": true, "singbox": true, "raw": true}
	if !validFormats[c.OutputFormat] {
		return fmt.Errorf("OUTPUT_FORMAT must be one of clash, singbox, raw; got %q", c.OutputFormat)
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.LogLevel] {
		return fmt.Errorf("LOG_LEVEL must be one of debug, info, warn, error; got %q", c.LogLevel)
	}

	return nil
}

// Addr returns the host:port string for the server to listen on.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
