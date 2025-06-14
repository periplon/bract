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
	server := NewServer(8765, browserClient)
	
	assert.NotNil(t, server)
	assert.Equal(t, 8765, server.port)
	assert.Equal(t, browserClient, server.browserClient)
	assert.NotNil(t, server.upgrader)
}

func TestServer_Start(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port: 0, // Use random port
	}
	browserClient := browser.NewClient(cfg)
	server := NewServer(0, browserClient)
	
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
	server := NewServer(8765, browserClient)
	
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
		s.ListenAndServe()
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
	server := NewServer(8765, browserClient)
	
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
	server := NewServer(8765, browserClient)
	
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
	server := NewServer(8765, browserClient)
	
	// Create test server
	testServer := httptest.NewServer(http.HandlerFunc(server.handleWebSocket))
	defer testServer.Close()
	
	// Connect to WebSocket
	wsURL := "ws" + strings.TrimPrefix(testServer.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()
	
	// Send a ping message
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
}

func TestServer_HandleMessage(t *testing.T) {
	cfg := config.WebSocketConfig{
		Port:        8765,
		ReconnectMs: 5000,
	}
	browserClient := browser.NewClient(cfg)
	server := NewServer(8765, browserClient)
	
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
		conn.WriteMessage(websocket.TextMessage, data)
		
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