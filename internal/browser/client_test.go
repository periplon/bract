package browser

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/periplon/bract/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConnection is a mock implementation of the Connection interface
type MockConnection struct {
	mock.Mock
}

func (m *MockConnection) SendCommand(action string, data interface{}) (string, error) {
	args := m.Called(action, data)
	return args.String(0), args.Error(1)
}

func TestNewClient(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port:        8765,
		ReconnectMs: 5000,
	}

	client := NewClient(cfg)

	assert.NotNil(t, client)
	assert.Equal(t, cfg, client.config)
	assert.Equal(t, -1, client.activeTabID)
	assert.Nil(t, client.connection)
}

func TestClient_SetConnection(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	mockConn := new(MockConnection)

	client.SetConnection(mockConn)

	assert.Equal(t, mockConn, client.connection)
}

func TestClient_RemoveConnection(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	mockConn := new(MockConnection)

	client.SetConnection(mockConn)
	client.RemoveConnection(mockConn)

	assert.Nil(t, client.connection)
}

func TestClient_RemoveConnection_DifferentConnection(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	mockConn1 := new(MockConnection)
	mockConn2 := new(MockConnection)

	client.SetConnection(mockConn1)
	client.RemoveConnection(mockConn2) // Try to remove different connection

	assert.Equal(t, mockConn1, client.connection) // Should not remove
}

func TestClient_WaitForConnection(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	mockConn := new(MockConnection)

	// Set connection after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		client.SetConnection(mockConn)
	}()

	ctx := context.Background()
	err := client.WaitForConnection(ctx, 1*time.Second)

	assert.NoError(t, err)
	assert.Equal(t, mockConn, client.connection)
}

func TestClient_WaitForConnection_Timeout(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})

	ctx := context.Background()
	err := client.WaitForConnection(ctx, 100*time.Millisecond)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout waiting for Chrome extension connection")
}

func TestClient_HandleResponse(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})

	// Store a pending response channel
	msgID := "test-123"
	respChan := make(chan Response, 1)
	client.pending.Store(msgID, respChan)

	// Handle response
	testData := json.RawMessage(`{"result": "success"}`)
	go client.HandleResponse(msgID, testData, "")

	// Check response was received
	select {
	case resp := <-respChan:
		assert.Equal(t, testData, resp.Data)
		assert.Empty(t, resp.Error)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	// Check pending was cleaned up
	_, exists := client.pending.Load(msgID)
	assert.False(t, exists)
}

func TestClient_HandleResponse_WithError(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})

	// Store a pending response channel
	msgID := "test-456"
	respChan := make(chan Response, 1)
	client.pending.Store(msgID, respChan)

	// Handle error response
	go client.HandleResponse(msgID, nil, "Test error")

	// Check error was received
	select {
	case resp := <-respChan:
		assert.Nil(t, resp.Data)
		assert.Equal(t, "Test error", resp.Error)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}
}

func TestClient_HandleEvent(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	client.activeTabID = 123

	// Handle tab closed event for active tab
	eventData := json.RawMessage(`{"tabId": 123}`)
	client.HandleEvent("tabClosed", eventData)

	assert.Equal(t, -1, client.activeTabID)
}

func TestClient_HandleEvent_DifferentTab(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	client.activeTabID = 123

	// Handle tab closed event for different tab
	eventData := json.RawMessage(`{"tabId": 456}`)
	client.HandleEvent("tabClosed", eventData)

	assert.Equal(t, 123, client.activeTabID) // Should not change
}

// Surfingkeys Integration Tests

func TestClient_ShowHints(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	expectedData := map[string]interface{}{
		"tabId":    123,
		"selector": "a",
		"action":   "click",
	}
	mockConn.On("SendCommand", "hints.show", expectedData).Return("msg-123", nil)

	// Simulate response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"hintsShown": 10}`)
		client.HandleResponse("msg-123", responseData, "")
	}()

	ctx := context.Background()
	result, err := client.ShowHints(ctx, 123, "a", "click")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockConn.AssertExpectations(t)
}

func TestClient_ClickHint(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	expectedData := map[string]interface{}{
		"tabId": 123,
		"index": 0,
		"text":  "Click me",
	}
	mockConn.On("SendCommand", "hints.click", expectedData).Return("msg-124", nil)

	// Simulate response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"clicked": true}`)
		client.HandleResponse("msg-124", responseData, "")
	}()

	ctx := context.Background()
	result, err := client.ClickHint(ctx, 123, "", 0, "Click me")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockConn.AssertExpectations(t)
}

func TestClient_Search(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	expectedData := map[string]interface{}{
		"query":  "test search",
		"engine": "google",
		"newTab": true,
	}
	mockConn.On("SendCommand", "search", expectedData).Return("msg-125", nil)

	// Simulate response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"searchInitiated": true}`)
		client.HandleResponse("msg-125", responseData, "")
	}()

	ctx := context.Background()
	result, err := client.Search(ctx, "test search", "google", true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockConn.AssertExpectations(t)
}

func TestClient_Find(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	expectedData := map[string]interface{}{
		"tabId":         123,
		"text":          "find me",
		"caseSensitive": true,
		"wholeWord":     false,
	}
	mockConn.On("SendCommand", "find", expectedData).Return("msg-126", nil)

	// Simulate response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"found": 3, "activeMatch": 1}`)
		client.HandleResponse("msg-126", responseData, "")
	}()

	ctx := context.Background()
	result, err := client.Find(ctx, 123, "find me", true, false)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockConn.AssertExpectations(t)
}

func TestClient_ReadClipboard(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	mockConn.On("SendCommand", "clipboard.read", nil).Return("msg-127", nil)

	// Store a pending response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"text": "clipboard content"}`)
		client.HandleResponse("msg-127", responseData, "")
	}()

	ctx := context.Background()
	result, err := client.ReadClipboard(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "clipboard content", result)
	mockConn.AssertExpectations(t)
}

func TestClient_WriteClipboard(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	expectedData := map[string]interface{}{
		"text":   "test content",
		"format": "text",
	}
	mockConn.On("SendCommand", "clipboard.write", expectedData).Return("msg-128", nil)

	// Simulate response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"success": true}`)
		client.HandleResponse("msg-128", responseData, "")
	}()

	ctx := context.Background()
	err := client.WriteClipboard(ctx, "test content", "text")

	assert.NoError(t, err)
	mockConn.AssertExpectations(t)
}

func TestClient_ShowOmnibar(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	expectedData := map[string]interface{}{
		"tabId": 123,
		"type":  "bookmarks",
		"query": "test",
	}
	mockConn.On("SendCommand", "omnibar.show", expectedData).Return("msg-129", nil)

	// Simulate response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"omnibarShown": true}`)
		client.HandleResponse("msg-129", responseData, "")
	}()

	ctx := context.Background()
	result, err := client.ShowOmnibar(ctx, 123, "bookmarks", "test")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockConn.AssertExpectations(t)
}

func TestClient_StartVisualMode(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	expectedData := map[string]interface{}{
		"tabId":         123,
		"selectElement": true,
	}
	mockConn.On("SendCommand", "visual.start", expectedData).Return("msg-130", nil)

	// Simulate response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"visualModeStarted": true}`)
		client.HandleResponse("msg-130", responseData, "")
	}()

	ctx := context.Background()
	result, err := client.StartVisualMode(ctx, 123, true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockConn.AssertExpectations(t)
}

func TestClient_GetPageTitle(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	expectedData := map[string]interface{}{
		"tabId": 123,
	}
	mockConn.On("SendCommand", "getPageTitle", expectedData).Return("msg-131", nil)

	// Store a pending response
	go func() {
		time.Sleep(10 * time.Millisecond)
		responseData := json.RawMessage(`{"title": "Test Page Title"}`)
		client.HandleResponse("msg-131", responseData, "")
	}()

	ctx := context.Background()
	result, err := client.GetPageTitle(ctx, 123)

	assert.NoError(t, err)
	assert.Equal(t, "Test Page Title", result)
	mockConn.AssertExpectations(t)
}

func TestClient_NoConnection(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})

	ctx := context.Background()
	_, err := client.ShowHints(ctx, 123, "a", "click")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no connection to Chrome extension")
}

func TestClient_Timeout(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	mockConn.On("SendCommand", "hints.show", mock.Anything).Return("msg-timeout", nil)

	// Don't send a response to trigger timeout
	ctx := context.Background()
	_, err := client.ShowHints(ctx, 123, "a", "click")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout waiting for Chrome extension response")
	mockConn.AssertExpectations(t)
}

func TestClient_ContextCancellation(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 5000})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	mockConn.On("SendCommand", "hints.show", mock.Anything).Return("msg-cancel", nil)

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := client.ShowHints(ctx, 123, "a", "click")

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	mockConn.AssertExpectations(t)
}

func TestClient_ChromeExtensionError(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := new(MockConnection)
	client.SetConnection(mockConn)

	mockConn.On("SendCommand", "hints.show", mock.Anything).Return("msg-error", nil)

	// Send error response
	go func() {
		time.Sleep(10 * time.Millisecond)
		client.HandleResponse("msg-error", nil, "Chrome extension error")
	}()

	ctx := context.Background()
	_, err := client.ShowHints(ctx, 123, "a", "click")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chrome extension error: Chrome extension error")
	mockConn.AssertExpectations(t)
}