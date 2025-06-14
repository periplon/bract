package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/periplon/bract/internal/browser"
	"github.com/periplon/bract/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port: 8765,
	}
	browserClient := browser.NewClient(cfg)
	allowedOrigins := []string{"http://localhost", "chrome-extension://*"}
	server := NewServer(8765, browserClient, allowedOrigins)

	assert.NotNil(t, server)
	assert.Equal(t, 8765, server.port)
	assert.Equal(t, browserClient, server.browserClient)
	assert.NotNil(t, server.upgrader)
	assert.Equal(t, allowedOrigins, server.allowedOrigins)
}

func TestServer_Start(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port: 0, // Use random port
	}
	browserClient := browser.NewClient(cfg)
	allowedOrigins := []string{"http://localhost", "chrome-extension://*"}
	server := NewServer(0, browserClient, allowedOrigins)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		err := server.Start(ctx)
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test health endpoint
	resp, err := http.Get("http://localhost:8765/health")
	if err == nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Cancel context to stop server
	cancel()

	// Wait for server to stop
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("Server did not stop within timeout")
	}
}

func TestServer_HandleWebSocket(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port: 8765,
	}
	browserClient := browser.NewClient(cfg)
	allowedOrigins := []string{"http://localhost", "chrome-extension://*"}
	server := NewServer(8765, browserClient, allowedOrigins)

	// Create a test request and response recorder
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Sec-WebSocket-Version", "13")

	// Use httptest to create a test server
	handler := http.HandlerFunc(server.handleWebSocket)
	s := &http.Server{Handler: handler}

	// Start test server
	go func() {
		_ = s.ListenAndServe()
	}()
	defer s.Close()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
}

func TestServer_MultipleConnections(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port: 8765,
	}
	browserClient := browser.NewClient(cfg)
	allowedOrigins := []string{"http://localhost", "chrome-extension://*"}
	server := NewServer(8765, browserClient, allowedOrigins)

	// Track connections
	var connCount int
	var mu sync.Mutex

	// Create test handler that counts connections
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := server.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		mu.Lock()
		connCount++
		mu.Unlock()

		// Keep connection open
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		connCount--
		mu.Unlock()
	})

	// Start test server
	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	// Create multiple connections
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")
	var conns []*websocket.Conn

	for i := 0; i < 3; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		conns = append(conns, ws)
	}

	// Give connections time to be established
	time.Sleep(50 * time.Millisecond)

	// Verify connection count
	mu.Lock()
	assert.Equal(t, 3, connCount)
	mu.Unlock()

	// Close connections
	for _, conn := range conns {
		conn.Close()
	}

	// Give time for cleanup
	time.Sleep(150 * time.Millisecond)

	// Verify all connections closed
	mu.Lock()
	assert.Equal(t, 0, connCount)
	mu.Unlock()
}

func TestServer_GetConnections(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port: 8765,
	}
	browserClient := browser.NewClient(cfg)
	allowedOrigins := []string{"http://localhost", "chrome-extension://*"}
	server := NewServer(8765, browserClient, allowedOrigins)

	// Initially no connections
	conns := server.GetConnections()
	assert.Empty(t, conns)

	// Add some test connections
	for i := 0; i < 3; i++ {
		conn := &Connection{
			ID: string(rune('A' + i)),
		}
		server.connections.Store(conn.ID, conn)
	}

	// Get connections
	conns = server.GetConnections()
	assert.Len(t, conns, 3)

	// Verify all connections are returned
	ids := make(map[string]bool)
	for _, conn := range conns {
		ids[conn.ID] = true
	}
	assert.True(t, ids["A"])
	assert.True(t, ids["B"])
	assert.True(t, ids["C"])
}

func TestServer_WebSocketIntegration(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port:        8765,
		ReconnectMs: 5000,
	}
	browserClient := browser.NewClient(cfg)
	allowedOrigins := []string{"http://localhost", "chrome-extension://*"}
	server := NewServer(8765, browserClient, allowedOrigins)

	// Create test server
	testServer := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
	defer testServer.Close()

	// Connect to WebSocket
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Test ping message
	pingMsg := Message{
		ID:   "test-ping",
		Type: "ping",
	}
	err = ws.WriteJSON(pingMsg)
	require.NoError(t, err)

	// Should receive pong response
	var response Message
	err = ws.ReadJSON(&response)
	require.NoError(t, err)
	assert.Equal(t, "test-ping", response.ID)
	assert.Equal(t, "pong", response.Type)

	// Test connected message
	connectedMsg := Message{
		ID:   "test-connected",
		Type: "connected",
	}
	err = ws.WriteJSON(connectedMsg)
	require.NoError(t, err)

	// Should receive ack response
	err = ws.ReadJSON(&response)
	require.NoError(t, err)
	assert.Equal(t, "test-connected", response.ID)
	assert.Equal(t, "ack", response.Type)
}

func TestServer_HandleMessage(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port:        8765,
		ReconnectMs: 5000,
	}
	browserClient := browser.NewClient(cfg)
	allowedOrigins := []string{"http://localhost", "chrome-extension://*"}
	server := NewServer(8765, browserClient, allowedOrigins)

	// Track responses
	responseChan := make(chan browser.Response, 1)

	// Create test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := server.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		connection := &Connection{
			ID:     "test-conn",
			conn:   conn,
			send:   make(chan []byte, 256),
			server: server,
		}

		// Register connection
		server.connections.Store(connection.ID, connection)
		browserClient.SetConnection(connection)

		// Start write pump
		go connection.writePump()

		// Send response message
		responseMsg := Message{
			ID:   "resp123",
			Type: "response",
			Data: json.RawMessage(`{"result":"success"}`),
		}

		data, _ := json.Marshal(responseMsg)
		_ = conn.WriteMessage(websocket.TextMessage, data)

		// Keep connection open
		time.Sleep(200 * time.Millisecond)
	}))
	defer testServer.Close()

	// Store response handler
	browserClient.HandleResponse("resp123", json.RawMessage(`{"result":"success"}`), "")

	// Connect to WebSocket
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Read and process messages
	go func() {
		for {
			var msg Message
			if err := ws.ReadJSON(&msg); err != nil {
				break
			}

			if msg.Type == "response" {
				responseChan <- browser.Response{
					Data:  msg.Data,
					Error: msg.Error,
				}
			}
		}
	}()

	// Wait for response
	select {
	case resp := <-responseChan:
		assert.Empty(t, resp.Error)
		assert.JSONEq(t, `{"result":"success"}`, string(resp.Data))
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for response")
	}
}

func TestServer_CheckOrigin(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigins []string
		origin         string
		expected       bool
	}{
		{
			name:           "Empty origin allowed",
			allowedOrigins: []string{"http://localhost"},
			origin:         "",
			expected:       true,
		},
		{
			name:           "Exact match localhost",
			allowedOrigins: []string{"http://localhost"},
			origin:         "http://localhost",
			expected:       true,
		},
		{
			name:           "Localhost with port",
			allowedOrigins: []string{"http://localhost"},
			origin:         "http://localhost:3000",
			expected:       true,
		},
		{
			name:           "Chrome extension wildcard",
			allowedOrigins: []string{"chrome-extension://*"},
			origin:         "chrome-extension://abcdefgh",
			expected:       true,
		},
		{
			name:           "Specific Chrome extension",
			allowedOrigins: []string{"chrome-extension://specific-id"},
			origin:         "chrome-extension://specific-id",
			expected:       true,
		},
		{
			name:           "Disallowed origin",
			allowedOrigins: []string{"http://localhost"},
			origin:         "http://example.com",
			expected:       false,
		},
		{
			name:           "HTTPS localhost",
			allowedOrigins: []string{"https://localhost"},
			origin:         "https://localhost:8443",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.WebSocketConfig{Port: 8765}
			browserClient := browser.NewClient(cfg)
			server := NewServer(8765, browserClient, tt.allowedOrigins)

			req := &http.Request{
				Header: http.Header{
					"Origin": []string{tt.origin},
				},
			}

			result := server.checkOrigin(req)
			assert.Equal(t, tt.expected, result)
		})
	}
}
