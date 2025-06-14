package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the server configuration
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	WebSocket WebSocketConfig `yaml:"websocket"`
	Browser   BrowserConfig   `yaml:"browser"`
	Logging   LoggingConfig   `yaml:"logging"`
}

// ServerConfig contains MCP server settings
type ServerConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// WebSocketConfig contains WebSocket server settings
type WebSocketConfig struct {
	Host           string   `yaml:"host"`
	Port           int      `yaml:"port"`
	ReconnectMs    int      `yaml:"reconnect_ms"`
	PingInterval   int      `yaml:"ping_interval"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// BrowserConfig contains browser automation settings
type BrowserConfig struct {
	DefaultTimeout int `yaml:"default_timeout"`
	MaxTabs        int `yaml:"max_tabs"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Name:    "Browser Automation Server",
			Version: "1.0.0",
		},
		WebSocket: WebSocketConfig{
			Host:         "localhost",
			Port:         8765,
			ReconnectMs:  5000,
			PingInterval: 30,
			AllowedOrigins: []string{
				"http://localhost",
				"https://localhost",
				"chrome-extension://*",
			},
		},
		Browser: BrowserConfig{
			DefaultTimeout: 30000,
			MaxTabs:        100,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

// Load loads configuration from file or returns default
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Try to load from file
	configPath := os.Getenv("MCP_BROWSER_CONFIG")
	if configPath == "" {
		// Try default locations
		for _, path := range []string{
			"./configs/config.yaml",
			"./config.yaml",
			filepath.Join(os.Getenv("HOME"), ".config/mcp-browser/config.yaml"),
		} {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
	}

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Override with environment variables if set
	if host := os.Getenv("MCP_BROWSER_WS_HOST"); host != "" {
		cfg.WebSocket.Host = host
	}
	if port := os.Getenv("MCP_BROWSER_WS_PORT"); port != "" {
		// Parse port from string
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			cfg.WebSocket.Port = p
		}
	}

	return cfg, nil
}
