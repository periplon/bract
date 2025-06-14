package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "Browser Automation Server", cfg.Server.Name)
	assert.Equal(t, "1.0.0", cfg.Server.Version)
	assert.Equal(t, "localhost", cfg.WebSocket.Host)
	assert.Equal(t, 8765, cfg.WebSocket.Port)
	assert.Equal(t, 5000, cfg.WebSocket.ReconnectMs)
	assert.Equal(t, 30, cfg.WebSocket.PingInterval)
	assert.Equal(t, 30000, cfg.Browser.DefaultTimeout)
	assert.Equal(t, 100, cfg.Browser.MaxTabs)
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

func TestLoad_DefaultWhenNoConfigFile(t *testing.T) {
	// Clear any existing config env var
	oldEnv := os.Getenv("MCP_BROWSER_CONFIG")
	os.Unsetenv("MCP_BROWSER_CONFIG")
	defer func() {
		if oldEnv != "" {
			os.Setenv("MCP_BROWSER_CONFIG", oldEnv)
		}
	}()

	// Create a temp directory that doesn't contain config files
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, DefaultConfig(), cfg)
}

func TestLoad_ValidYAMLFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  name: "Test Server"
  version: "2.0.0"
websocket:
  host: "0.0.0.0"
  port: 9876
  reconnect_ms: 10000
  ping_interval: 60
browser:
  default_timeout: 60000
  max_tabs: 50
logging:
  level: "debug"
  format: "text"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set env var to point to test config
	oldEnv := os.Getenv("MCP_BROWSER_CONFIG")
	os.Setenv("MCP_BROWSER_CONFIG", configPath)
	defer func() {
		if oldEnv != "" {
			os.Setenv("MCP_BROWSER_CONFIG", oldEnv)
		} else {
			os.Unsetenv("MCP_BROWSER_CONFIG")
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "Test Server", cfg.Server.Name)
	assert.Equal(t, "2.0.0", cfg.Server.Version)
	assert.Equal(t, "0.0.0.0", cfg.WebSocket.Host)
	assert.Equal(t, 9876, cfg.WebSocket.Port)
	assert.Equal(t, 10000, cfg.WebSocket.ReconnectMs)
	assert.Equal(t, 60, cfg.WebSocket.PingInterval)
	assert.Equal(t, 60000, cfg.Browser.DefaultTimeout)
	assert.Equal(t, 50, cfg.Browser.MaxTabs)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "text", cfg.Logging.Format)
}

func TestLoad_PartialYAMLFile(t *testing.T) {
	// Create temporary config file with partial config
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
websocket:
  port: 8080
browser:
  max_tabs: 200
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set env var to point to test config
	oldEnv := os.Getenv("MCP_BROWSER_CONFIG")
	os.Setenv("MCP_BROWSER_CONFIG", configPath)
	defer func() {
		if oldEnv != "" {
			os.Setenv("MCP_BROWSER_CONFIG", oldEnv)
		} else {
			os.Unsetenv("MCP_BROWSER_CONFIG")
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)

	// Modified values
	assert.Equal(t, 8080, cfg.WebSocket.Port)
	assert.Equal(t, 200, cfg.Browser.MaxTabs)

	// Default values
	assert.Equal(t, "Browser Automation Server", cfg.Server.Name)
	assert.Equal(t, "1.0.0", cfg.Server.Version)
	assert.Equal(t, "localhost", cfg.WebSocket.Host)
	assert.Equal(t, 5000, cfg.WebSocket.ReconnectMs)
	assert.Equal(t, 30, cfg.WebSocket.PingInterval)
	assert.Equal(t, 30000, cfg.Browser.DefaultTimeout)
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

func TestLoad_InvalidYAMLFile(t *testing.T) {
	// Create temporary config file with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  port: "not a number"
  host: 123 # should be string
invalid yaml here
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set env var to point to test config
	oldEnv := os.Getenv("MCP_BROWSER_CONFIG")
	os.Setenv("MCP_BROWSER_CONFIG", configPath)
	defer func() {
		if oldEnv != "" {
			os.Setenv("MCP_BROWSER_CONFIG", oldEnv)
		} else {
			os.Unsetenv("MCP_BROWSER_CONFIG")
		}
	}()

	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoad_EmptyYAMLFile(t *testing.T) {
	// Create empty config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(""), 0644)
	require.NoError(t, err)

	// Set env var to point to test config
	oldEnv := os.Getenv("MCP_BROWSER_CONFIG")
	os.Setenv("MCP_BROWSER_CONFIG", configPath)
	defer func() {
		if oldEnv != "" {
			os.Setenv("MCP_BROWSER_CONFIG", oldEnv)
		} else {
			os.Unsetenv("MCP_BROWSER_CONFIG")
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)

	// Should return default config
	assert.Equal(t, DefaultConfig(), cfg)
}

func TestLoad_EnvironmentVariableOverrides(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
websocket:
  host: "localhost"
  port: 8765
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Save old env vars
	oldConfig := os.Getenv("MCP_BROWSER_CONFIG")
	oldHost := os.Getenv("MCP_BROWSER_WS_HOST")
	oldPort := os.Getenv("MCP_BROWSER_WS_PORT")

	// Set environment variables
	os.Setenv("MCP_BROWSER_CONFIG", configPath)
	os.Setenv("MCP_BROWSER_WS_HOST", "0.0.0.0")
	os.Setenv("MCP_BROWSER_WS_PORT", "9999")

	defer func() {
		// Restore old env vars
		if oldConfig != "" {
			os.Setenv("MCP_BROWSER_CONFIG", oldConfig)
		} else {
			os.Unsetenv("MCP_BROWSER_CONFIG")
		}
		if oldHost != "" {
			os.Setenv("MCP_BROWSER_WS_HOST", oldHost)
		} else {
			os.Unsetenv("MCP_BROWSER_WS_HOST")
		}
		if oldPort != "" {
			os.Setenv("MCP_BROWSER_WS_PORT", oldPort)
		} else {
			os.Unsetenv("MCP_BROWSER_WS_PORT")
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)

	// Environment variables should override file config
	assert.Equal(t, "0.0.0.0", cfg.WebSocket.Host)
	assert.Equal(t, 9999, cfg.WebSocket.Port)
}

func TestLoad_InvalidPortEnvironmentVariable(t *testing.T) {
	// Save old env vars
	oldPort := os.Getenv("MCP_BROWSER_WS_PORT")

	// Set invalid port
	os.Setenv("MCP_BROWSER_WS_PORT", "not-a-number")

	defer func() {
		if oldPort != "" {
			os.Setenv("MCP_BROWSER_WS_PORT", oldPort)
		} else {
			os.Unsetenv("MCP_BROWSER_WS_PORT")
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)

	// Should use default port when env var is invalid
	assert.Equal(t, 8765, cfg.WebSocket.Port)
}

func TestLoad_DefaultLocations(t *testing.T) {
	// Clear config env var
	oldEnv := os.Getenv("MCP_BROWSER_CONFIG")
	os.Unsetenv("MCP_BROWSER_CONFIG")
	defer func() {
		if oldEnv != "" {
			os.Setenv("MCP_BROWSER_CONFIG", oldEnv)
		}
	}()

	// Create config in current directory
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
websocket:
  port: 7777
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, 7777, cfg.WebSocket.Port)
}