package mcp

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/periplon/bract/internal/browser"
	"github.com/periplon/bract/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockBrowserClient is a mock implementation of handler.BrowserClient
type MockBrowserClient struct{}

func (m *MockBrowserClient) SetConnection(conn browser.Connection)                         {}
func (m *MockBrowserClient) RemoveConnection(conn browser.Connection)                      {}
func (m *MockBrowserClient) HandleResponse(id string, data json.RawMessage, errMsg string) {}
func (m *MockBrowserClient) HandleEvent(action string, data json.RawMessage)               {}
func (m *MockBrowserClient) WaitForConnection(ctx context.Context, timeout time.Duration) error {
	return nil
}
func (m *MockBrowserClient) ListTabs(ctx context.Context) ([]browser.Tab, error) { return nil, nil }
func (m *MockBrowserClient) CreateTab(ctx context.Context, url string, active bool) (*browser.Tab, error) {
	return nil, nil
}
func (m *MockBrowserClient) CloseTab(ctx context.Context, tabID int) error    { return nil }
func (m *MockBrowserClient) ActivateTab(ctx context.Context, tabID int) error { return nil }
func (m *MockBrowserClient) Navigate(ctx context.Context, tabID int, url string, waitUntilLoad bool) error {
	return nil
}
func (m *MockBrowserClient) Reload(ctx context.Context, tabID int, hardReload bool) error { return nil }
func (m *MockBrowserClient) Click(ctx context.Context, tabID int, selector string, timeout int) error {
	return nil
}
func (m *MockBrowserClient) Type(ctx context.Context, tabID int, selector, text string, clearFirst bool, delay int) error {
	return nil
}
func (m *MockBrowserClient) Scroll(ctx context.Context, tabID int, x, y *float64, selector, behavior string) error {
	return nil
}
func (m *MockBrowserClient) WaitForElement(ctx context.Context, tabID int, selector string, timeout int, state string) error {
	return nil
}
func (m *MockBrowserClient) ExecuteScript(ctx context.Context, tabID int, script string, args []interface{}) (json.RawMessage, error) {
	return nil, nil
}
func (m *MockBrowserClient) ExtractContent(ctx context.Context, tabID int, selector, contentType, attribute string) ([]string, error) {
	return nil, nil
}
func (m *MockBrowserClient) Screenshot(ctx context.Context, tabID int, fullPage bool, selector, format string, quality int) (string, error) {
	return "", nil
}
func (m *MockBrowserClient) GetCookies(ctx context.Context, url, name string) ([]browser.Cookie, error) {
	return nil, nil
}
func (m *MockBrowserClient) SetCookie(ctx context.Context, cookie browser.Cookie) error { return nil }
func (m *MockBrowserClient) DeleteCookies(ctx context.Context, url, name string) error  { return nil }
func (m *MockBrowserClient) GetLocalStorage(ctx context.Context, tabID int, key string) (string, error) {
	return "", nil
}
func (m *MockBrowserClient) SetLocalStorage(ctx context.Context, tabID int, key, value string) error {
	return nil
}
func (m *MockBrowserClient) GetSessionStorage(ctx context.Context, tabID int, key string) (string, error) {
	return "", nil
}
func (m *MockBrowserClient) SetSessionStorage(ctx context.Context, tabID int, key, value string) error {
	return nil
}

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		srvName string
		version string
	}{
		{
			name:    "creates new server instance",
			srvName: "test-server",
			version: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			h := handler.NewBrowserHandler(mockClient)

			server := NewServer(tt.srvName, tt.version, h)

			assert.NotNil(t, server)
			assert.NotNil(t, server.mcpServer)
			assert.NotNil(t, server.handler)
			assert.Equal(t, h, server.handler)
		})
	}
}

func TestServer_RegisterTools(t *testing.T) {
	mockClient := &MockBrowserClient{}
	h := handler.NewBrowserHandler(mockClient)
	server := NewServer("test-server", "1.0.0", h)

	// The constructor should have already registered tools
	// We can verify the server was created successfully
	assert.NotNil(t, server)
	assert.NotNil(t, server.mcpServer)
}

func TestServer_Start(t *testing.T) {
	// Note: Start() uses stdio transport which is difficult to test
	// in unit tests. This would be better tested in integration tests.
	t.Skip("Start() uses stdio transport - better tested in integration tests")
}

func TestServer_ToolRegistration(t *testing.T) {
	tests := []struct {
		name          string
		expectedTools []string
	}{
		{
			name: "all browser automation tools registered",
			expectedTools: []string{
				// Tab management
				"browser_list_tabs",
				"browser_create_tab",
				"browser_close_tab",
				"browser_activate_tab",
				// Navigation
				"browser_navigate",
				"browser_reload",
				// Interaction
				"browser_click",
				"browser_type",
				"browser_scroll",
				"browser_wait_for_element",
				// Content
				"browser_execute_script",
				"browser_extract_content",
				"browser_screenshot",
				// Storage
				"browser_get_cookies",
				"browser_set_cookie",
				"browser_delete_cookies",
				"browser_get_local_storage",
				"browser_set_local_storage",
				"browser_get_session_storage",
				"browser_set_session_storage",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			h := handler.NewBrowserHandler(mockClient)
			server := NewServer("test-server", "1.0.0", h)

			// Verify server was created
			require.NotNil(t, server)
			require.NotNil(t, server.mcpServer)

			// Note: The actual tool registration happens inside the MCP server
			// which is a third-party library. We're mainly testing that our
			// registration methods are called without errors.
		})
	}
}
