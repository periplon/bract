package handler

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/periplon/bract/internal/browser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBrowserClient is a mock implementation of BrowserClient interface
type MockBrowserClient struct {
	mock.Mock
}

func (m *MockBrowserClient) SetConnection(conn browser.Connection) {
	m.Called(conn)
}

func (m *MockBrowserClient) RemoveConnection(conn browser.Connection) {
	m.Called(conn)
}

func (m *MockBrowserClient) HandleResponse(id string, data json.RawMessage, errMsg string) {
	m.Called(id, data, errMsg)
}

func (m *MockBrowserClient) HandleEvent(action string, data json.RawMessage) {
	m.Called(action, data)
}

func (m *MockBrowserClient) WaitForConnection(ctx context.Context, timeout time.Duration) error {
	args := m.Called(ctx, timeout)
	return args.Error(0)
}

func (m *MockBrowserClient) ListTabs(ctx context.Context) ([]browser.Tab, error) {
	args := m.Called(ctx)
	return args.Get(0).([]browser.Tab), args.Error(1)
}

func (m *MockBrowserClient) CreateTab(ctx context.Context, url string, active bool) (*browser.Tab, error) {
	args := m.Called(ctx, url, active)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*browser.Tab), args.Error(1)
}

func (m *MockBrowserClient) CloseTab(ctx context.Context, tabID int) error {
	args := m.Called(ctx, tabID)
	return args.Error(0)
}

func (m *MockBrowserClient) ActivateTab(ctx context.Context, tabID int) error {
	args := m.Called(ctx, tabID)
	return args.Error(0)
}

func (m *MockBrowserClient) Navigate(ctx context.Context, tabID int, url string, waitUntilLoad bool) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, url, waitUntilLoad)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) Reload(ctx context.Context, tabID int, hardReload bool) error {
	args := m.Called(ctx, tabID, hardReload)
	return args.Error(0)
}

func (m *MockBrowserClient) Click(ctx context.Context, tabID int, selector string, timeout int) error {
	args := m.Called(ctx, tabID, selector, timeout)
	return args.Error(0)
}

func (m *MockBrowserClient) Type(ctx context.Context, tabID int, selector, text string, clearFirst bool, delay int) error {
	args := m.Called(ctx, tabID, selector, text, clearFirst, delay)
	return args.Error(0)
}

func (m *MockBrowserClient) Scroll(ctx context.Context, tabID int, x, y *float64, selector, behavior string) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, x, y, selector, behavior)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) WaitForElement(ctx context.Context, tabID int, selector string, timeout int, state string) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, selector, timeout, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) ExecuteScript(ctx context.Context, tabID int, script string, args []interface{}) (json.RawMessage, error) {
	callArgs := m.Called(ctx, tabID, script, args)
	if callArgs.Get(0) == nil {
		return nil, callArgs.Error(1)
	}
	return callArgs.Get(0).(json.RawMessage), callArgs.Error(1)
}

func (m *MockBrowserClient) ExtractContent(ctx context.Context, tabID int, selector, contentType, attribute string) ([]string, error) {
	args := m.Called(ctx, tabID, selector, contentType, attribute)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockBrowserClient) ExtractText(ctx context.Context, tabID int, selector string) (string, error) {
	args := m.Called(ctx, tabID, selector)
	return args.String(0), args.Error(1)
}

func (m *MockBrowserClient) Screenshot(ctx context.Context, tabID int, fullPage bool, selector, format string, quality int) (string, error) {
	args := m.Called(ctx, tabID, fullPage, selector, format, quality)
	return args.String(0), args.Error(1)
}

func (m *MockBrowserClient) GetCookies(ctx context.Context, url, name string) ([]browser.Cookie, error) {
	args := m.Called(ctx, url, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]browser.Cookie), args.Error(1)
}

func (m *MockBrowserClient) SetCookie(ctx context.Context, cookie browser.Cookie) (json.RawMessage, error) {
	args := m.Called(ctx, cookie)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) DeleteCookies(ctx context.Context, url, name string) error {
	args := m.Called(ctx, url, name)
	return args.Error(0)
}

func (m *MockBrowserClient) GetLocalStorage(ctx context.Context, tabID int, key string) (string, error) {
	args := m.Called(ctx, tabID, key)
	return args.String(0), args.Error(1)
}

func (m *MockBrowserClient) SetLocalStorage(ctx context.Context, tabID int, key, value string) error {
	args := m.Called(ctx, tabID, key, value)
	return args.Error(0)
}

func (m *MockBrowserClient) ClearLocalStorage(ctx context.Context, tabID int) error {
	args := m.Called(ctx, tabID)
	return args.Error(0)
}

func (m *MockBrowserClient) GetSessionStorage(ctx context.Context, tabID int, key string) (string, error) {
	args := m.Called(ctx, tabID, key)
	return args.String(0), args.Error(1)
}

func (m *MockBrowserClient) SetSessionStorage(ctx context.Context, tabID int, key, value string) error {
	args := m.Called(ctx, tabID, key, value)
	return args.Error(0)
}

func (m *MockBrowserClient) ClearSessionStorage(ctx context.Context, tabID int) error {
	args := m.Called(ctx, tabID)
	return args.Error(0)
}

func (m *MockBrowserClient) GetActionables(ctx context.Context, tabID int) ([]browser.Actionable, error) {
	args := m.Called(ctx, tabID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]browser.Actionable), args.Error(1)
}

func (m *MockBrowserClient) GetAccessibilitySnapshot(ctx context.Context, tabID int, interestingOnly bool, root string) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, interestingOnly, root)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

// Helper function to extract text from mcp.Content
func getTextFromContent(t *testing.T, content mcp.Content) string {
	// Try both pointer and value types since the interface could contain either
	switch tc := content.(type) {
	case *mcp.TextContent:
		return tc.Text
	case mcp.TextContent:
		return tc.Text
	default:
		t.Fatalf("Expected TextContent type, got %T", content)
		return ""
	}
}

func TestNewBrowserHandler(t *testing.T) {
	mockClient := &MockBrowserClient{}
	handler := NewBrowserHandler(mockClient)

	assert.NotNil(t, handler)
	assert.Equal(t, mockClient, handler.client)
}

func TestBrowserHandler_WaitForConnection(t *testing.T) {
	tests := []struct {
		name        string
		request     mcp.CallToolRequest
		setupMock   func(*MockBrowserClient)
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "successful connection",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_wait_for_connection",
					Arguments: map[string]interface{}{},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("WaitForConnection", mock.Anything, 30*time.Second).Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, "Successfully connected to browser extension", text)
			},
		},
		{
			name: "connection with custom timeout",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_wait_for_connection",
					Arguments: map[string]interface{}{
						"timeout": 10,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("WaitForConnection", mock.Anything, 10*time.Second).Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, "Successfully connected to browser extension", text)
			},
		},
		{
			name: "connection timeout",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_wait_for_connection",
					Arguments: map[string]interface{}{},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("WaitForConnection", mock.Anything, 30*time.Second).Return(errors.New("connection timeout"))
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to connect to browser")
				assert.Contains(t, text, "connection timeout")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			handler := NewBrowserHandler(mockClient)

			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			result, err := handler.WaitForConnection(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestBrowserHandler_ListTabs(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockBrowserClient)
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "list tabs successfully",
			setupMock: func(m *MockBrowserClient) {
				tabs := []browser.Tab{
					{ID: 1, URL: "https://example.com", Title: "Example", Active: true},
					{ID: 2, URL: "https://google.com", Title: "Google", Active: false},
				}
				m.On("ListTabs", mock.Anything).Return(tabs, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])

				// Should return JSON array
				var tabs []browser.Tab
				err := json.Unmarshal([]byte(text), &tabs)
				require.NoError(t, err)
				assert.Len(t, tabs, 2)
				assert.Equal(t, "https://example.com", tabs[0].URL)
				assert.True(t, tabs[0].Active)
			},
		},
		{
			name: "list tabs with error",
			setupMock: func(m *MockBrowserClient) {
				m.On("ListTabs", mock.Anything).Return([]browser.Tab(nil), errors.New("browser not connected"))
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to list tabs")
				assert.Contains(t, text, "browser not connected")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			handler := NewBrowserHandler(mockClient)

			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_list_tabs",
					Arguments: map[string]interface{}{},
				},
			}

			result, err := handler.ListTabs(context.Background(), request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestBrowserHandler_CreateTab(t *testing.T) {
	tests := []struct {
		name        string
		request     mcp.CallToolRequest
		setupMock   func(*MockBrowserClient)
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "create tab successfully",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_create_tab",
					Arguments: map[string]interface{}{
						"url":    "https://example.com",
						"active": true,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				tab := &browser.Tab{
					ID:     3,
					URL:    "https://example.com",
					Title:  "Example",
					Active: true,
				}
				m.On("CreateTab", mock.Anything, "https://example.com", true).Return(tab, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])

				var tab browser.Tab
				err := json.Unmarshal([]byte(text), &tab)
				require.NoError(t, err)
				assert.Equal(t, 3, tab.ID)
				assert.Equal(t, "https://example.com", tab.URL)
				assert.True(t, tab.Active)
			},
		},
		{
			name: "create tab with defaults",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_create_tab",
					Arguments: map[string]interface{}{},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				tab := &browser.Tab{
					ID:     4,
					URL:    "about:blank",
					Title:  "New Tab",
					Active: true,
				}
				m.On("CreateTab", mock.Anything, "about:blank", true).Return(tab, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])

				var tab browser.Tab
				err := json.Unmarshal([]byte(text), &tab)
				require.NoError(t, err)
				assert.Equal(t, "about:blank", tab.URL)
			},
		},
		{
			name: "create tab with error",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_create_tab",
					Arguments: map[string]interface{}{
						"url": "https://example.com",
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("CreateTab", mock.Anything, "https://example.com", true).Return((*browser.Tab)(nil), errors.New("failed to create tab"))
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to create tab")
				assert.Contains(t, text, "failed to create tab")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			handler := NewBrowserHandler(mockClient)

			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			result, err := handler.CreateTab(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestBrowserHandler_CloseTab(t *testing.T) {
	tests := []struct {
		name        string
		request     mcp.CallToolRequest
		setupMock   func(*MockBrowserClient)
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "close tab successfully",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_close_tab",
					Arguments: map[string]interface{}{
						"tabId": 5,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("CloseTab", mock.Anything, 5).Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, "Closed tab 5", text)
			},
		},
		{
			name: "close tab missing id",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_close_tab",
					Arguments: map[string]interface{}{},
				},
			},
			setupMock: nil,
			wantErr:   false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "required argument \"tabId\" not found")
			},
		},
		{
			name: "close tab with error",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_close_tab",
					Arguments: map[string]interface{}{
						"tabId": 6,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("CloseTab", mock.Anything, 6).Return(errors.New("tab not found"))
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to close tab")
				assert.Contains(t, text, "tab not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			handler := NewBrowserHandler(mockClient)

			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			result, err := handler.CloseTab(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestBrowserHandler_Navigate(t *testing.T) {
	tests := []struct {
		name        string
		request     mcp.CallToolRequest
		setupMock   func(*MockBrowserClient)
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "navigate successfully",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_navigate",
					Arguments: map[string]interface{}{
						"url":           "https://example.com",
						"waitUntilLoad": true,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("Navigate", mock.Anything, 0, "https://example.com", true).Return(json.RawMessage(`{"success": true}`), nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, `{"success": true}`, text)
			},
		},
		{
			name: "navigate with specific tab",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_navigate",
					Arguments: map[string]interface{}{
						"url":   "https://google.com",
						"tabId": 7,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("Navigate", mock.Anything, 7, "https://google.com", true).Return(json.RawMessage(`{"success": true}`), nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, `{"success": true}`, text)
			},
		},
		{
			name: "navigate missing url",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_navigate",
					Arguments: map[string]interface{}{},
				},
			},
			setupMock: nil,
			wantErr:   false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "required argument \"url\" not found")
			},
		},
		{
			name: "navigate with error",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_navigate",
					Arguments: map[string]interface{}{
						"url": "https://invalid.site",
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("Navigate", mock.Anything, 0, "https://invalid.site", true).Return(json.RawMessage(nil), errors.New("failed to navigate"))
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to navigate")
				assert.Contains(t, text, "failed to navigate")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			handler := NewBrowserHandler(mockClient)

			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			result, err := handler.Navigate(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestBrowserHandler_Click(t *testing.T) {
	tests := []struct {
		name        string
		request     mcp.CallToolRequest
		setupMock   func(*MockBrowserClient)
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "click successfully",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_click",
					Arguments: map[string]interface{}{
						"selector": "#submit-button",
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("Click", mock.Anything, 0, "#submit-button", 30000).Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, "Clicked on element: #submit-button", text)
			},
		},
		{
			name: "click with custom timeout",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_click",
					Arguments: map[string]interface{}{
						"selector": ".btn",
						"timeout":  5000,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("Click", mock.Anything, 0, ".btn", 5000).Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, "Clicked on element: .btn", text)
			},
		},
		{
			name: "click missing selector",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_click",
					Arguments: map[string]interface{}{},
				},
			},
			setupMock: nil,
			wantErr:   false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "required argument \"selector\" not found")
			},
		},
		{
			name: "click element not found",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_click",
					Arguments: map[string]interface{}{
						"selector": "#missing",
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("Click", mock.Anything, 0, "#missing", 30000).Return(errors.New("element not found"))
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to click")
				assert.Contains(t, text, "element not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			handler := NewBrowserHandler(mockClient)

			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			result, err := handler.Click(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestBrowserHandler_GetActionables(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockBrowserClient)
		request     mcp.CallToolRequest
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "get actionables successfully",
			setupMock: func(m *MockBrowserClient) {
				actionables := []browser.Actionable{
					{
						LabelNumber: 0,
						Description: "Submit",
						Type:        "button",
						Selector:    "#submit-btn",
					},
					{
						LabelNumber: 1,
						Description: "Email input",
						Type:        "input",
						Selector:    "input[type=\"email\"]",
					},
					{
						LabelNumber: 2,
						Description: "Home",
						Type:        "link",
						Selector:    "a[href=\"/\"]",
					},
				}
				m.On("GetActionables", mock.Anything, 0).Return(actionables, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_get_actionables",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])

				// Should return JSON array
				var actionablesResult []browser.Actionable
				err := json.Unmarshal([]byte(text), &actionablesResult)
				require.NoError(t, err)
				assert.Len(t, actionablesResult, 3)

				assert.Equal(t, 0, actionablesResult[0].LabelNumber)
				assert.Equal(t, "Submit", actionablesResult[0].Description)
				assert.Equal(t, "button", actionablesResult[0].Type)
				assert.Equal(t, "#submit-btn", actionablesResult[0].Selector)

				assert.Equal(t, 1, actionablesResult[1].LabelNumber)
				assert.Equal(t, "Email input", actionablesResult[1].Description)
				assert.Equal(t, "input", actionablesResult[1].Type)
			},
		},
		{
			name: "get actionables with specific tab ID",
			setupMock: func(m *MockBrowserClient) {
				actionables := []browser.Actionable{
					{
						LabelNumber: 0,
						Description: "Click me",
						Type:        "button",
						Selector:    "button.action",
					},
				}
				m.On("GetActionables", mock.Anything, 123).Return(actionables, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_get_actionables",
					Arguments: map[string]interface{}{
						"tabId": 123,
					},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])

				var actionablesResult []browser.Actionable
				err := json.Unmarshal([]byte(text), &actionablesResult)
				require.NoError(t, err)
				assert.Len(t, actionablesResult, 1)

				assert.Equal(t, 0, actionablesResult[0].LabelNumber)
				assert.Equal(t, "Click me", actionablesResult[0].Description)
				assert.Equal(t, "button", actionablesResult[0].Type)
				assert.Equal(t, "button.action", actionablesResult[0].Selector)
			},
		},
		{
			name: "get actionables empty page",
			setupMock: func(m *MockBrowserClient) {
				m.On("GetActionables", mock.Anything, 0).Return([]browser.Actionable{}, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_get_actionables",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])

				var actionablesResult []browser.Actionable
				err := json.Unmarshal([]byte(text), &actionablesResult)
				require.NoError(t, err)
				assert.Len(t, actionablesResult, 0)
			},
		},
		{
			name: "get actionables with error",
			setupMock: func(m *MockBrowserClient) {
				m.On("GetActionables", mock.Anything, 0).Return([]browser.Actionable(nil), errors.New("page not loaded"))
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_get_actionables",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to get actionables")
				assert.Contains(t, text, "page not loaded")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			handler := NewBrowserHandler(mockClient)

			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			result, err := handler.GetActionables(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestBrowserHandler_GetAccessibilitySnapshot(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockBrowserClient)
		request     mcp.CallToolRequest
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "get accessibility snapshot successfully",
			setupMock: func(m *MockBrowserClient) {
				snapshot := json.RawMessage(`{
					"snapshot": {
						"role": "RootWebArea",
						"name": "Test Page",
						"children": [
							{
								"role": "heading",
								"name": "Welcome",
								"level": 1
							},
							{
								"role": "button",
								"name": "Submit",
								"disabled": false
							}
						]
					}
				}`)
				m.On("GetAccessibilitySnapshot", mock.Anything, 0, true, "").Return(snapshot, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_get_accessibility_snapshot",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])

				// Should return valid JSON
				var snapshot map[string]interface{}
				err := json.Unmarshal([]byte(text), &snapshot)
				require.NoError(t, err)
				assert.Equal(t, "RootWebArea", snapshot["role"])
				assert.Equal(t, "Test Page", snapshot["name"])

				children := snapshot["children"].([]interface{})
				assert.Len(t, children, 2)
			},
		},
		{
			name: "get accessibility snapshot with specific parameters",
			setupMock: func(m *MockBrowserClient) {
				snapshot := json.RawMessage(`{
					"snapshot": {
						"role": "region",
						"name": "Main Content"
					}
				}`)
				m.On("GetAccessibilitySnapshot", mock.Anything, 123, false, "#main").Return(snapshot, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_get_accessibility_snapshot",
					Arguments: map[string]interface{}{
						"tabId":           123,
						"interestingOnly": false,
						"root":            "#main",
					},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])

				var snapshot map[string]interface{}
				err := json.Unmarshal([]byte(text), &snapshot)
				require.NoError(t, err)
				assert.Equal(t, "region", snapshot["role"])
				assert.Equal(t, "Main Content", snapshot["name"])
			},
		},
		{
			name: "get accessibility snapshot with error",
			setupMock: func(m *MockBrowserClient) {
				m.On("GetAccessibilitySnapshot", mock.Anything, 0, true, "").Return(json.RawMessage(nil), errors.New("page not accessible"))
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_get_accessibility_snapshot",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to get accessibility snapshot")
				assert.Contains(t, text, "page not accessible")
			},
		},
		{
			name: "get accessibility snapshot with null result",
			setupMock: func(m *MockBrowserClient) {
				snapshot := json.RawMessage(`{"snapshot": null}`)
				m.On("GetAccessibilitySnapshot", mock.Anything, 0, true, "").Return(snapshot, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_get_accessibility_snapshot",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				// Should return empty object when snapshot is null
				assert.Equal(t, "{}", text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockBrowserClient{}
			handler := NewBrowserHandler(mockClient)

			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			result, err := handler.GetAccessibilitySnapshot(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}
