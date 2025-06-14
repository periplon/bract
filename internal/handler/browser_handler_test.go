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
	mockArgs := m.Called(ctx, tabID, script, args)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(json.RawMessage), mockArgs.Error(1)
}

func (m *MockBrowserClient) ExtractContent(ctx context.Context, tabID int, selector, contentType, attribute string) ([]string, error) {
	args := m.Called(ctx, tabID, selector, contentType, attribute)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
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

func (m *MockBrowserClient) GetSessionStorage(ctx context.Context, tabID int, key string) (string, error) {
	args := m.Called(ctx, tabID, key)
	return args.String(0), args.Error(1)
}

func (m *MockBrowserClient) SetSessionStorage(ctx context.Context, tabID int, key, value string) error {
	args := m.Called(ctx, tabID, key, value)
	return args.Error(0)
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
			name: "wait for connection successfully",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_wait_for_connection",
					Arguments: map[string]interface{}{
						"timeout": 5.0,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("WaitForConnection", mock.Anything, 5*time.Second).Return(nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Successfully connected to browser extension")
			},
		},
		{
			name: "wait for connection with default timeout",
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
				assert.Contains(t, text, "Successfully connected to browser extension")
			},
		},
		{
			name: "wait for connection timeout",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_wait_for_connection",
					Arguments: map[string]interface{}{
						"timeout": 2.0,
					},
				},
			},
			setupMock: func(m *MockBrowserClient) {
				m.On("WaitForConnection", mock.Anything, 2*time.Second).Return(errors.New("timeout waiting for Chrome extension connection"))
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to connect to browser")
				assert.Contains(t, text, "timeout")
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
					{ID: 1, URL: "https://example.com", Title: "Example", Index: 0, Active: true},
					{ID: 2, URL: "https://google.com", Title: "Google", Index: 1, Active: false},
				}
				m.On("ListTabs", mock.Anything).Return(tabs, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				
				// Should return JSON array
				var tabsResult []browser.Tab
				err := json.Unmarshal([]byte(text), &tabsResult)
				require.NoError(t, err)
				assert.Len(t, tabsResult, 2)
				assert.Equal(t, 1, tabsResult[0].ID)
				assert.Equal(t, "Example", tabsResult[0].Title)
				assert.True(t, tabsResult[0].Active)
			},
		},
		{
			name: "list tabs with error",
			setupMock: func(m *MockBrowserClient) {
				m.On("ListTabs", mock.Anything).Return([]browser.Tab{}, errors.New("connection failed"))
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "Failed to list tabs")
				assert.Contains(t, text, "connection failed")
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
					Name:      "list_tabs",
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
		setupMock   func(*MockBrowserClient)
		request     mcp.CallToolRequest
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "create tab successfully",
			setupMock: func(m *MockBrowserClient) {
				tab := &browser.Tab{
					ID:    3,
					URL:   "https://example.com",
					Title: "Example",
					Index: 2,
				}
				m.On("CreateTab", mock.Anything, "https://example.com", true).Return(tab, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "create_tab",
					Arguments: map[string]interface{}{
						"url":    "https://example.com",
						"active": true,
					},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				// Check that the result is valid JSON
				var tab browser.Tab
				err := json.Unmarshal([]byte(text), &tab)
				assert.NoError(t, err)
				assert.Equal(t, 3, tab.ID)
				assert.Equal(t, "https://example.com", tab.URL)
				assert.Equal(t, "Example", tab.Title)
			},
		},
		{
			name: "create tab without URL",
			setupMock: func(m *MockBrowserClient) {
				tab := &browser.Tab{
					ID:    4,
					URL:   "about:blank",
					Title: "New Tab",
					Index: 3,
				}
				m.On("CreateTab", mock.Anything, "about:blank", true).Return(tab, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "create_tab",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				// Check that the result is valid JSON
				var tab browser.Tab
				err := json.Unmarshal([]byte(text), &tab)
				assert.NoError(t, err)
				assert.Equal(t, 4, tab.ID)
				assert.Equal(t, "about:blank", tab.URL)
				assert.Equal(t, "New Tab", tab.Title)
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

func TestBrowserHandler_Screenshot(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockBrowserClient)
		request     mcp.CallToolRequest
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "take screenshot successfully",
			setupMock: func(m *MockBrowserClient) {
				m.On("Screenshot", mock.Anything, 0, false, "", "png", 90).Return("data:image/png;base64,iVBORw0KGgoAAAANS", nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "screenshot",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 2) // Text + Image

				// First content should be text
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, "Screenshot captured", text)

				// Second content should be image
				var imageContent mcp.ImageContent
				switch ic := result.Content[1].(type) {
				case *mcp.ImageContent:
					imageContent = *ic
				case mcp.ImageContent:
					imageContent = ic
				default:
					t.Fatalf("Expected ImageContent type, got %T", result.Content[1])
				}
				assert.Equal(t, "image", imageContent.Type)
				assert.Equal(t, "image/png", imageContent.MIMEType)
				assert.Equal(t, "iVBORw0KGgoAAAANS", imageContent.Data)
			},
		},
		{
			name: "take full page screenshot",
			setupMock: func(m *MockBrowserClient) {
				m.On("Screenshot", mock.Anything, 1, true, "", "jpeg", 80).Return("data:image/jpeg;base64,/9j/4AAQSkZJRg", nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "screenshot",
					Arguments: map[string]interface{}{
						"tabId":    1,
						"fullPage": true,
						"format":   "jpeg",
						"quality":  80,
					},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 2) // Text + Image

				// First content should be text
				text := getTextFromContent(t, result.Content[0])
				assert.Equal(t, "Screenshot captured", text)

				// Second content should be image
				var imageContent mcp.ImageContent
				switch ic := result.Content[1].(type) {
				case *mcp.ImageContent:
					imageContent = *ic
				case mcp.ImageContent:
					imageContent = ic
				default:
					t.Fatalf("Expected ImageContent type, got %T", result.Content[1])
				}
				assert.Equal(t, "image/jpeg", imageContent.MIMEType)
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

			result, err := handler.Screenshot(context.Background(), tt.request)

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

func TestBrowserHandler_ExecuteScript(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockBrowserClient)
		request     mcp.CallToolRequest
		wantErr     bool
		checkResult func(*testing.T, *mcp.CallToolResult)
	}{
		{
			name: "execute script with string result",
			setupMock: func(m *MockBrowserClient) {
				result := json.RawMessage(`"Page Title"`)
				m.On("ExecuteScript", mock.Anything, 0, "return document.title", []interface{}(nil)).Return(result, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "execute_script",
					Arguments: map[string]interface{}{
						"script": "return document.title",
					},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				// Should return raw JSON result
				assert.Equal(t, `"Page Title"`, text)
			},
		},
		{
			name: "execute script with args",
			setupMock: func(m *MockBrowserClient) {
				result := json.RawMessage(`3`)
				m.On("ExecuteScript", mock.Anything, 0, "return arguments[0] + arguments[1]", []interface{}{1.0, 2.0}).Return(result, nil)
			},
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "execute_script",
					Arguments: map[string]interface{}{
						"script": "return arguments[0] + arguments[1]",
						"args":   []interface{}{1.0, 2.0},
					},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				// Should return raw JSON result
				assert.Equal(t, "3", text)
			},
		},
		{
			name: "execute script missing script parameter",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "execute_script",
					Arguments: map[string]interface{}{},
				},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *mcp.CallToolResult) {
				assert.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text := getTextFromContent(t, result.Content[0])
				assert.Contains(t, text, "required")
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

			result, err := handler.ExecuteScript(context.Background(), tt.request)

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
