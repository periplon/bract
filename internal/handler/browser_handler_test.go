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
)

// MockBrowserClient is a mock implementation of BrowserClient interface
type MockBrowserClient struct {
	mock.Mock
}

// Connection management
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

// Surfingkeys MCP Integration
func (m *MockBrowserClient) ShowHints(ctx context.Context, tabID int, selector, action string) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, selector, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) ClickHint(ctx context.Context, tabID int, selector string, index int, text string) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, selector, index, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) Search(ctx context.Context, query, engine string, newTab bool) (json.RawMessage, error) {
	args := m.Called(ctx, query, engine, newTab)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) Find(ctx context.Context, tabID int, text string, caseSensitive, wholeWord bool) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, text, caseSensitive, wholeWord)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) ReadClipboard(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockBrowserClient) WriteClipboard(ctx context.Context, text, format string) error {
	args := m.Called(ctx, text, format)
	return args.Error(0)
}

func (m *MockBrowserClient) ShowOmnibar(ctx context.Context, tabID int, barType, query string) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, barType, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) StartVisualMode(ctx context.Context, tabID int, selectElement bool) (json.RawMessage, error) {
	args := m.Called(ctx, tabID, selectElement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockBrowserClient) GetPageTitle(ctx context.Context, tabID int) (string, error) {
	args := m.Called(ctx, tabID)
	return args.String(0), args.Error(1)
}

// Tests

func TestNewBrowserHandler(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	assert.NotNil(t, handler)
	assert.Equal(t, mockClient, handler.client)
}

func TestBrowserHandler_WaitForConnection(t *testing.T) {
	tests := []struct {
		name           string
		timeout        float64
		clientError    error
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Success",
			timeout:        30,
			clientError:    nil,
			expectedResult: "Successfully connected to browser extension",
			expectError:    false,
		},
		{
			name:           "Connection Failed",
			timeout:        10,
			clientError:    errors.New("connection timeout"),
			expectedResult: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockBrowserClient)
			handler := NewBrowserHandler(mockClient)

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "browser_wait_for_connection",
					Arguments: map[string]interface{}{
						"timeout": tt.timeout,
					},
				},
			}

			mockClient.On("WaitForConnection", mock.Anything, time.Duration(tt.timeout)*time.Second).
				Return(tt.clientError)

			result, err := handler.WaitForConnection(context.Background(), request)

			assert.NoError(t, err) // Handler should not return error, but ToolResult

			if tt.expectError {
				assert.True(t, result.IsError)
				textContent, ok := result.Content[0].(mcp.TextContent)
				assert.True(t, ok)
				assert.Contains(t, textContent.Text, "Failed to connect to browser")
			} else {
				assert.False(t, result.IsError)
				textContent, ok := result.Content[0].(mcp.TextContent)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedResult, textContent.Text)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestBrowserHandler_ShowHints(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_hints_show",
			Arguments: map[string]interface{}{
				"selector": "a",
				"action":   "click",
				"tabId":    123,
			},
		},
	}

	expectedResponse := json.RawMessage(`{"hintsShown": 10}`)
	mockClient.On("ShowHints", mock.Anything, 123, "a", "click").
		Return(expectedResponse, nil)

	result, err := handler.ShowHints(context.Background(), request)

	assert.NoError(t, err)
	assert.False(t, result.IsError)
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Equal(t, string(expectedResponse), textContent.Text)
	mockClient.AssertExpectations(t)
}

func TestBrowserHandler_ClickHint(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "Click by selector",
			params: map[string]interface{}{
				"selector": "a.link",
				"tabId":    123,
			},
			expectError: false,
		},
		{
			name: "Click by index",
			params: map[string]interface{}{
				"index": 5,
				"tabId": 123,
			},
			expectError: false,
		},
		{
			name: "Click by text",
			params: map[string]interface{}{
				"text":  "Click me",
				"tabId": 123,
			},
			expectError: false,
		},
		{
			name:        "No selector, index, or text",
			params:      map[string]interface{}{"tabId": 123},
			expectError: true,
			errorMsg:    "Must provide either selector, index, or text parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockBrowserClient)
			handler := NewBrowserHandler(mockClient)

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "browser_hints_click",
					Arguments: tt.params,
				},
			}

			if !tt.expectError {
				expectedResponse := json.RawMessage(`{"clicked": true}`)
				selector := ""
				if s, ok := tt.params["selector"]; ok {
					selector = s.(string)
				}
				index := -1
				if i, ok := tt.params["index"]; ok {
					index = i.(int)
				}
				text := ""
				if t, ok := tt.params["text"]; ok {
					text = t.(string)
				}

				mockClient.On("ClickHint", mock.Anything, 123, selector, index, text).
					Return(expectedResponse, nil)
			}

			result, err := handler.ClickHint(context.Background(), request)

			assert.NoError(t, err) // Handler should not return error, but ToolResult

			if tt.expectError {
				assert.True(t, result.IsError)
				textContent, ok := result.Content[0].(mcp.TextContent)
			assert.True(t, ok)
			assert.Contains(t, textContent.Text, tt.errorMsg)
			} else {
				assert.False(t, result.IsError)
				mockClient.AssertExpectations(t)
			}
		})
	}
}

func TestBrowserHandler_Search(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_search",
			Arguments: map[string]interface{}{
				"query":  "test search",
				"engine": "google",
				"newTab": true,
			},
		},
	}

	expectedResponse := json.RawMessage(`{"searchInitiated": true}`)
	mockClient.On("Search", mock.Anything, "test search", "google", true).
		Return(expectedResponse, nil)

	result, err := handler.Search(context.Background(), request)

	assert.NoError(t, err)
	assert.False(t, result.IsError)
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Equal(t, string(expectedResponse), textContent.Text)
	mockClient.AssertExpectations(t)
}

func TestBrowserHandler_Find(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_find",
			Arguments: map[string]interface{}{
				"text":          "find me",
				"caseSensitive": true,
				"wholeWord":     false,
				"tabId":         123,
			},
		},
	}

	expectedResponse := json.RawMessage(`{"found": 3, "activeMatch": 1}`)
	mockClient.On("Find", mock.Anything, 123, "find me", true, false).
		Return(expectedResponse, nil)

	result, err := handler.Find(context.Background(), request)

	assert.NoError(t, err)
	assert.False(t, result.IsError)
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Equal(t, string(expectedResponse), textContent.Text)
	mockClient.AssertExpectations(t)
}

func TestBrowserHandler_ReadClipboard(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "browser_clipboard_read",
			Arguments: map[string]interface{}{},
		},
	}

	mockClient.On("ReadClipboard", mock.Anything).
		Return("clipboard content", nil)

	result, err := handler.ReadClipboard(context.Background(), request)

	assert.NoError(t, err)
	assert.False(t, result.IsError)
	
	// Parse the JSON response
	var response map[string]string
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	err = json.Unmarshal([]byte(textContent.Text), &response)
	assert.NoError(t, err)
	assert.Equal(t, "clipboard content", response["text"])
	
	mockClient.AssertExpectations(t)
}

func TestBrowserHandler_WriteClipboard(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_clipboard_write",
			Arguments: map[string]interface{}{
				"text":   "test content",
				"format": "text",
			},
		},
	}

	mockClient.On("WriteClipboard", mock.Anything, "test content", "text").
		Return(nil)

	result, err := handler.WriteClipboard(context.Background(), request)

	assert.NoError(t, err)
	assert.False(t, result.IsError)
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Equal(t, "Wrote to clipboard: test content", textContent.Text)
	mockClient.AssertExpectations(t)
}

func TestBrowserHandler_ShowOmnibar(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_omnibar",
			Arguments: map[string]interface{}{
				"type":  "bookmarks",
				"query": "test",
				"tabId": 123,
			},
		},
	}

	expectedResponse := json.RawMessage(`{"omnibarShown": true}`)
	mockClient.On("ShowOmnibar", mock.Anything, 123, "bookmarks", "test").
		Return(expectedResponse, nil)

	result, err := handler.ShowOmnibar(context.Background(), request)

	assert.NoError(t, err)
	assert.False(t, result.IsError)
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Equal(t, string(expectedResponse), textContent.Text)
	mockClient.AssertExpectations(t)
}

func TestBrowserHandler_StartVisualMode(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_visual_mode",
			Arguments: map[string]interface{}{
				"selectElement": true,
				"tabId":         123,
			},
		},
	}

	expectedResponse := json.RawMessage(`{"visualModeStarted": true}`)
	mockClient.On("StartVisualMode", mock.Anything, 123, true).
		Return(expectedResponse, nil)

	result, err := handler.StartVisualMode(context.Background(), request)

	assert.NoError(t, err)
	assert.False(t, result.IsError)
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Equal(t, string(expectedResponse), textContent.Text)
	mockClient.AssertExpectations(t)
}

func TestBrowserHandler_GetPageTitle(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_get_page_title",
			Arguments: map[string]interface{}{
				"tabId": 123,
			},
		},
	}

	mockClient.On("GetPageTitle", mock.Anything, 123).
		Return("Test Page Title", nil)

	result, err := handler.GetPageTitle(context.Background(), request)

	assert.NoError(t, err)
	assert.False(t, result.IsError)
	
	// Parse the JSON response
	var response map[string]string
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	err = json.Unmarshal([]byte(textContent.Text), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Page Title", response["title"])
	
	mockClient.AssertExpectations(t)
}

// Test error cases

func TestBrowserHandler_ShowHints_Error(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_hints_show",
			Arguments: map[string]interface{}{
				"selector": "a",
				"tabId":    123,
			},
		},
	}

	mockClient.On("ShowHints", mock.Anything, 123, "a", "").
		Return(json.RawMessage(nil), errors.New("extension error"))

	result, err := handler.ShowHints(context.Background(), request)

	assert.NoError(t, err) // Handler should not return error, but ToolResult
	assert.True(t, result.IsError)
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Contains(t, textContent.Text, "Failed to show hints")
	mockClient.AssertExpectations(t)
}

func TestBrowserHandler_Search_MissingQuery(t *testing.T) {
	mockClient := new(MockBrowserClient)
	handler := NewBrowserHandler(mockClient)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "browser_search",
			Arguments: map[string]interface{}{
				"engine": "google",
			},
		},
	}

	result, err := handler.Search(context.Background(), request)

	assert.NoError(t, err) // Handler should not return error, but ToolResult
	assert.True(t, result.IsError)
	textContent, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Contains(t, textContent.Text, "query")
}