package browser

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/periplon/bract/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockConnection is a mock implementation of Connection interface
type MockConnection struct {
	mock.Mock
}

func (m *MockConnection) SendCommand(action string, data interface{}) (string, error) {
	args := m.Called(action, data)
	return args.String(0), args.Error(1)
}

func TestNewClient(t *testing.T) {
	cfg := config.WebSocketConfig{
		Host:         "localhost",
		Port:         8765,
		ReconnectMs:  5000,
		PingInterval: 30,
	}

	client := NewClient(cfg)
	assert.NotNil(t, client)
	assert.Equal(t, cfg, client.config)
	assert.Equal(t, -1, client.activeTabID)
}

func TestClient_SetConnection(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	mockConn := &MockConnection{}

	client.SetConnection(mockConn)

	client.mu.RLock()
	assert.Equal(t, mockConn, client.connection)
	client.mu.RUnlock()
}

func TestClient_RemoveConnection(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	mockConn := &MockConnection{}

	client.SetConnection(mockConn)
	client.mu.RLock()
	assert.NotNil(t, client.connection)
	client.mu.RUnlock()

	client.RemoveConnection(mockConn)
	client.mu.RLock()
	assert.Nil(t, client.connection)
	client.mu.RUnlock()
}

func TestClient_RemoveConnection_DifferentConnection(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	mockConn1 := &MockConnection{}
	mockConn2 := &MockConnection{}

	client.SetConnection(mockConn1)
	client.RemoveConnection(mockConn2) // Try to remove different connection

	client.mu.RLock()
	assert.Equal(t, mockConn1, client.connection) // Original connection should remain
	client.mu.RUnlock()
}

func TestClient_HandleResponse(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})

	// Create a channel to receive response
	respChan := make(chan Response, 1)
	client.pending.Store("test-id", respChan)

	// Handle response
	testData := json.RawMessage(`{"result": "success"}`)
	client.HandleResponse("test-id", testData, "")

	// Verify response was received
	select {
	case resp := <-respChan:
		assert.Equal(t, testData, resp.Data)
		assert.Empty(t, resp.Error)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Response not received")
	}

	// Verify pending entry was removed
	_, exists := client.pending.Load("test-id")
	assert.False(t, exists)
}

func TestClient_HandleResponse_WithError(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})

	// Create a channel to receive response
	respChan := make(chan Response, 1)
	client.pending.Store("test-id", respChan)

	// Handle response with error
	client.HandleResponse("test-id", nil, "test error")

	// Verify error response was received
	select {
	case resp := <-respChan:
		assert.Nil(t, resp.Data)
		assert.Equal(t, "test error", resp.Error)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Response not received")
	}
}

func TestClient_HandleEvent(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	client.activeTabID = 123

	// Handle tab closed event
	eventData := json.RawMessage(`{"tabId": 123}`)
	client.HandleEvent("tabClosed", eventData)

	// Verify active tab was cleared
	assert.Equal(t, -1, client.activeTabID)
}

func TestClient_HandleEvent_DifferentTab(t *testing.T) {
	client := NewClient(config.WebSocketConfig{})
	client.activeTabID = 123

	// Handle tab closed event for different tab
	eventData := json.RawMessage(`{"tabId": 456}`)
	client.HandleEvent("tabClosed", eventData)

	// Verify active tab was not changed
	assert.Equal(t, 123, client.activeTabID)
}

func TestClient_ListTabs(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockConnection)
		hasConn      bool
		expectedTabs []Tab
		wantErr      bool
		errMsg       string
	}{
		{
			name: "list tabs successfully",
			setupMock: func(m *MockConnection) {
				m.On("SendCommand", "listTabs", nil).Return("msg-123", nil)
			},
			hasConn: true,
			expectedTabs: []Tab{
				{ID: 1, URL: "https://example.com", Title: "Example", Index: 0, Active: true},
				{ID: 2, URL: "https://google.com", Title: "Google", Index: 1, Active: false},
			},
			wantErr: false,
		},
		{
			name:    "list tabs without connection",
			hasConn: false,
			wantErr: true,
			errMsg:  "no connection",
		},
		{
			name: "list tabs send error",
			setupMock: func(m *MockConnection) {
				m.On("SendCommand", "listTabs", nil).Return("", errors.New("send failed"))
			},
			hasConn: true,
			wantErr: true,
			errMsg:  "send failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
			mockConn := &MockConnection{}

			if tt.hasConn {
				client.SetConnection(mockConn)
				if tt.setupMock != nil {
					tt.setupMock(mockConn)
				}

				// If successful, simulate response
				if !tt.wantErr && tt.expectedTabs != nil {
					go func() {
						time.Sleep(10 * time.Millisecond)
						tabsJSON, _ := json.Marshal(tt.expectedTabs)
						client.HandleResponse("msg-123", tabsJSON, "")
					}()
				}
			}

			tabs, err := client.ListTabs(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTabs, tabs)
			}

			if tt.hasConn {
				mockConn.AssertExpectations(t)
			}
		})
	}
}

func TestClient_CreateTab(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		active      bool
		setupMock   func(*MockConnection)
		hasConn     bool
		expectedTab *Tab
		wantErr     bool
		errMsg      string
	}{
		{
			name:   "create tab successfully",
			url:    "https://example.com",
			active: true,
			setupMock: func(m *MockConnection) {
				params := map[string]interface{}{
					"url":    "https://example.com",
					"active": true,
				}
				m.On("SendCommand", "createTab", params).Return("msg-456", nil)
			},
			hasConn: true,
			expectedTab: &Tab{
				ID:    3,
				URL:   "https://example.com",
				Title: "Example",
				Index: 2,
			},
			wantErr: false,
		},
		{
			name:    "create tab without connection",
			url:     "https://example.com",
			active:  false,
			hasConn: false,
			wantErr: true,
			errMsg:  "no connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
			mockConn := &MockConnection{}

			if tt.hasConn {
				client.SetConnection(mockConn)
				if tt.setupMock != nil {
					tt.setupMock(mockConn)
				}

				// If successful, simulate response
				if !tt.wantErr && tt.expectedTab != nil {
					go func() {
						time.Sleep(10 * time.Millisecond)
						tabJSON, _ := json.Marshal(tt.expectedTab)
						client.HandleResponse("msg-456", tabJSON, "")
					}()
				}
			}

			tab, err := client.CreateTab(context.Background(), tt.url, tt.active)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTab, tab)
				if tt.active {
					assert.Equal(t, tt.expectedTab.ID, client.activeTabID)
				}
			}

			if tt.hasConn {
				mockConn.AssertExpectations(t)
			}
		})
	}
}

func TestClient_CloseTab(t *testing.T) {
	tests := []struct {
		name      string
		tabID     int
		setupMock func(*MockConnection)
		hasConn   bool
		wantErr   bool
		errMsg    string
	}{
		{
			name:  "close tab successfully",
			tabID: 123,
			setupMock: func(m *MockConnection) {
				params := map[string]interface{}{
					"tabId": 123,
				}
				m.On("SendCommand", "closeTab", params).Return("msg-789", nil)
			},
			hasConn: true,
			wantErr: false,
		},
		{
			name:    "close tab without connection",
			tabID:   123,
			hasConn: false,
			wantErr: true,
			errMsg:  "no connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
			mockConn := &MockConnection{}

			if tt.hasConn {
				client.SetConnection(mockConn)
				if tt.setupMock != nil {
					tt.setupMock(mockConn)
				}

				// If successful, simulate response
				if !tt.wantErr {
					go func() {
						time.Sleep(10 * time.Millisecond)
						client.HandleResponse("msg-789", json.RawMessage(`{}`), "")
					}()
				}
			}

			err := client.CloseTab(context.Background(), tt.tabID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
			}

			if tt.hasConn {
				mockConn.AssertExpectations(t)
			}
		})
	}
}

func TestClient_Navigate(t *testing.T) {
	tests := []struct {
		name          string
		tabID         int
		url           string
		waitUntilLoad bool
		setupMock     func(*MockConnection)
		hasConn       bool
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "navigate successfully",
			tabID:         123,
			url:           "https://example.com",
			waitUntilLoad: true,
			setupMock: func(m *MockConnection) {
				params := map[string]interface{}{
					"tabId":         123,
					"url":           "https://example.com",
					"waitUntilLoad": true,
				}
				m.On("SendCommand", "navigate", params).Return("msg-nav", nil)
			},
			hasConn: true,
			wantErr: false,
		},
		{
			name:          "navigate with active tab",
			tabID:         0,
			url:           "https://example.com",
			waitUntilLoad: false,
			setupMock: func(m *MockConnection) {
				params := map[string]interface{}{
					"tabId":         456, // activeTabID
					"url":           "https://example.com",
					"waitUntilLoad": false,
				}
				m.On("SendCommand", "navigate", params).Return("msg-nav2", nil)
			},
			hasConn: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
			client.activeTabID = 456
			mockConn := &MockConnection{}

			if tt.hasConn {
				client.SetConnection(mockConn)
				if tt.setupMock != nil {
					tt.setupMock(mockConn)
				}

				// Simulate response
				if !tt.wantErr {
					go func() {
						time.Sleep(10 * time.Millisecond)
						msgID := "msg-nav"
						if tt.tabID == 0 {
							msgID = "msg-nav2"
						}
						client.HandleResponse(msgID, json.RawMessage(`{}`), "")
					}()
				}
			}

			_, err := client.Navigate(context.Background(), tt.tabID, tt.url, tt.waitUntilLoad)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
			}

			if tt.hasConn {
				mockConn.AssertExpectations(t)
			}
		})
	}
}

func TestClient_ExecuteScript(t *testing.T) {
	tests := []struct {
		name           string
		tabID          int
		script         string
		args           []interface{}
		setupMock      func(*MockConnection)
		hasConn        bool
		expectedResult json.RawMessage
		wantErr        bool
		errMsg         string
	}{
		{
			name:   "execute script with result",
			tabID:  123,
			script: "return document.title",
			args:   nil,
			setupMock: func(m *MockConnection) {
				params := map[string]interface{}{
					"tabId":  123,
					"script": "return document.title",
					"args":   []interface{}(nil),
				}
				m.On("SendCommand", "executeScript", params).Return("msg-exec", nil)
			},
			hasConn:        true,
			expectedResult: json.RawMessage(`"Page Title"`),
			wantErr:        false,
		},
		{
			name:   "execute script with args",
			tabID:  0, // Use active tab
			script: "return arguments[0] + arguments[1]",
			args:   []interface{}{1, 2},
			setupMock: func(m *MockConnection) {
				params := map[string]interface{}{
					"tabId":  456, // activeTabID
					"script": "return arguments[0] + arguments[1]",
					"args":   []interface{}{1, 2},
				}
				m.On("SendCommand", "executeScript", params).Return("msg-exec2", nil)
			},
			hasConn:        true,
			expectedResult: json.RawMessage(`3`),
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
			client.activeTabID = 456
			mockConn := &MockConnection{}

			if tt.hasConn {
				client.SetConnection(mockConn)
				if tt.setupMock != nil {
					tt.setupMock(mockConn)
				}

				// Simulate response
				if !tt.wantErr {
					go func() {
						time.Sleep(10 * time.Millisecond)
						msgID := "msg-exec"
						if tt.tabID == 0 {
							msgID = "msg-exec2"
						}
						client.HandleResponse(msgID, tt.expectedResult, "")
					}()
				}
			}

			result, err := client.ExecuteScript(context.Background(), tt.tabID, tt.script, tt.args)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			if tt.hasConn {
				mockConn.AssertExpectations(t)
			}
		})
	}
}

func TestClient_Screenshot(t *testing.T) {
	tests := []struct {
		name            string
		tabID           int
		fullPage        bool
		selector        string
		format          string
		quality         int
		setupMock       func(*MockConnection)
		hasConn         bool
		expectedDataURL string
		wantErr         bool
		errMsg          string
	}{
		{
			name:     "take screenshot successfully",
			tabID:    123,
			fullPage: false,
			selector: "",
			format:   "png",
			quality:  90,
			setupMock: func(m *MockConnection) {
				params := map[string]interface{}{
					"tabId":    123,
					"fullPage": false,
					"format":   "png",
					"quality":  90,
				}
				m.On("SendCommand", "screenshot", params).Return("msg-screen", nil)
			},
			hasConn:         true,
			expectedDataURL: "data:image/png;base64,iVBORw0KGgoAAAANS",
			wantErr:         false,
		},
		{
			name:     "take full page screenshot with selector",
			tabID:    0,
			fullPage: true,
			selector: "#content",
			format:   "jpeg",
			quality:  80,
			setupMock: func(m *MockConnection) {
				params := map[string]interface{}{
					"tabId":    456, // activeTabID
					"fullPage": true,
					"format":   "jpeg",
					"quality":  80,
					"selector": "#content",
				}
				m.On("SendCommand", "screenshot", params).Return("msg-screen2", nil)
			},
			hasConn:         true,
			expectedDataURL: "data:image/jpeg;base64,/9j/4AAQSkZJRg",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
			client.activeTabID = 456
			mockConn := &MockConnection{}

			if tt.hasConn {
				client.SetConnection(mockConn)
				if tt.setupMock != nil {
					tt.setupMock(mockConn)
				}

				// Simulate response
				if !tt.wantErr {
					go func() {
						time.Sleep(10 * time.Millisecond)
						msgID := "msg-screen"
						if tt.tabID == 0 {
							msgID = "msg-screen2"
						}
						responseData := json.RawMessage(`{"dataUrl":"` + tt.expectedDataURL + `"}`)
						client.HandleResponse(msgID, responseData, "")
					}()
				}
			}

			dataURL, err := client.Screenshot(context.Background(), tt.tabID, tt.fullPage, tt.selector, tt.format, tt.quality)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedDataURL, dataURL)
			}

			if tt.hasConn {
				mockConn.AssertExpectations(t)
			}
		})
	}
}

func TestClient_Timeout(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 50}) // Short timeout
	mockConn := &MockConnection{}
	client.SetConnection(mockConn)

	// Setup mock to return message ID but never send response
	mockConn.On("SendCommand", "listTabs", nil).Return("msg-timeout", nil)

	// Try to list tabs (will timeout)
	_, err := client.ListTabs(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")

	mockConn.AssertExpectations(t)
}

func TestClient_ContextCancellation(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 5000})
	mockConn := &MockConnection{}
	client.SetConnection(mockConn)

	// Setup mock to return message ID but never send response
	mockConn.On("SendCommand", "listTabs", nil).Return("msg-cancel", nil)

	// Create a context that we'll cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Start the operation in a goroutine
	errCh := make(chan error)
	go func() {
		_, err := client.ListTabs(ctx)
		errCh <- err
	}()

	// Cancel the context after a short delay
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Wait for the error
	err := <-errCh
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)

	mockConn.AssertExpectations(t)
}

func TestClient_ChromeExtensionError(t *testing.T) {
	client := NewClient(config.WebSocketConfig{ReconnectMs: 100})
	mockConn := &MockConnection{}
	client.SetConnection(mockConn)

	// Setup mock
	mockConn.On("SendCommand", "listTabs", nil).Return("msg-error", nil)

	// Simulate error response
	go func() {
		time.Sleep(10 * time.Millisecond)
		client.HandleResponse("msg-error", nil, "Chrome extension error: tabs permission denied")
	}()

	// Try to list tabs
	_, err := client.ListTabs(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "chrome extension error")
	assert.Contains(t, err.Error(), "tabs permission denied")

	mockConn.AssertExpectations(t)
}
