package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/periplon/bract/internal/browser"
	"github.com/periplon/bract/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnection_SendCommand(t *testing.T) {
	// Create test server
	msgChan := make(chan *Message, 10)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		// Read messages and store them
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}
			msgChan <- &msg

			// Send response
			response := Message{
				ID:   msg.ID,
				Type: "response",
				Data: json.RawMessage(`{"status":"ok"}`),
			}
			if err := conn.WriteJSON(response); err != nil {
				break
			}
		}
	}))
	defer server.Close()

	// Create connection
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	conn := &Connection{
		ID:   uuid.New().String(),
		conn: ws,
		send: make(chan []byte, 256),
	}

	// Start write pump
	go conn.writePump()

	// Test sending command
	msgID, err := conn.SendCommand("testAction", map[string]string{"key": "value"})
	require.NoError(t, err)
	assert.NotEmpty(t, msgID)

	// Verify message was received
	select {
	case msg := <-msgChan:
		assert.Equal(t, msgID, msg.ID)
		assert.Equal(t, "command", msg.Type)
		assert.Equal(t, "testAction", msg.Action)

		var data map[string]string
		err := json.Unmarshal(msg.Data, &data)
		require.NoError(t, err)
		assert.Equal(t, "value", data["key"])
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

func TestConnection_SendMessage(t *testing.T) {
	// Create test server
	msgChan := make(chan *Message, 10)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Read messages
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}
			msgChan <- &msg
		}
	}))
	defer server.Close()

	// Create connection
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	conn := &Connection{
		ID:   uuid.New().String(),
		conn: ws,
		send: make(chan []byte, 256),
	}

	// Start write pump
	go conn.writePump()

	// Test sending message
	msg := &Message{
		ID:     uuid.New().String(),
		Type:   "test",
		Action: "testAction",
		Data:   json.RawMessage(`{"test":"data"}`),
	}

	err = conn.SendMessage(msg)
	require.NoError(t, err)

	// Verify message was received
	select {
	case received := <-msgChan:
		assert.Equal(t, msg.ID, received.ID)
		assert.Equal(t, msg.Type, received.Type)
		assert.Equal(t, msg.Action, received.Action)
		assert.JSONEq(t, string(msg.Data), string(received.Data))
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

func TestConnection_WritePump(t *testing.T) {
	// Create test server that counts messages
	var msgCount int
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Count received messages
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
			mu.Lock()
			msgCount++
			mu.Unlock()
		}
	}))
	defer server.Close()

	// Create connection
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	conn := &Connection{
		ID:   uuid.New().String(),
		conn: ws,
		send: make(chan []byte, 256),
	}

	// Start write pump
	go conn.writePump()

	// Send multiple messages through the channel
	for i := 0; i < 5; i++ {
		msg := map[string]interface{}{
			"id":   i,
			"data": "test",
		}
		data, err := json.Marshal(msg)
		require.NoError(t, err)

		select {
		case conn.send <- data:
		case <-time.After(time.Second):
			t.Fatal("Timeout sending message")
		}
	}

	// Wait for messages to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify all messages were sent
	mu.Lock()
	assert.Equal(t, 5, msgCount)
	mu.Unlock()
}

func TestConnection_ReadPump(t *testing.T) {
	// Create channels to track messages
	pongChan := make(chan *Message, 10)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Send different message types
		messages := []Message{
			{
				ID:   "resp1",
				Type: "response",
				Data: json.RawMessage(`{"result":"success"}`),
			},
			{
				ID:     "event1",
				Type:   "event",
				Action: "tabClosed",
				Data:   json.RawMessage(`{"tabId":123}`),
			},
			{
				ID:   "ping1",
				Type: "ping",
			},
		}

		for _, msg := range messages {
			if err := conn.WriteJSON(msg); err != nil {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}

		// Read pong response
		go func() {
			for {
				var msg Message
				if err := conn.ReadJSON(&msg); err != nil {
					break
				}
				if msg.Type == "pong" {
					pongChan <- &msg
				}
			}
		}()

		// Keep connection open for pong response
		time.Sleep(200 * time.Millisecond)
	}))
	defer server.Close()

	// Create browser client with config
	cfg := config.WebSocketConfig{
		Port: 8765,
	}
	browserClient := browser.NewClient(cfg)

	// Create websocket server
	wsServer := &Server{
		browserClient: browserClient,
		connections:   sync.Map{},
	}

	// Create connection
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	conn := &Connection{
		ID:     uuid.New().String(),
		conn:   ws,
		send:   make(chan []byte, 256),
		server: wsServer,
	}

	// Register connection
	wsServer.connections.Store(conn.ID, conn)
	browserClient.SetConnection(conn)

	// Start write pump for pong responses
	go conn.writePump()

	// Start read pump
	go conn.readPump()

	// Wait for messages to be processed
	time.Sleep(300 * time.Millisecond)

	// Verify pong was sent
	select {
	case pong := <-pongChan:
		assert.Equal(t, "ping1", pong.ID)
		assert.Equal(t, "pong", pong.Type)
	case <-time.After(time.Second):
		t.Fatal("Did not receive pong response")
	}
}

func TestConnection_SendBufferFull(t *testing.T) {
	conn := &Connection{
		ID:   uuid.New().String(),
		send: make(chan []byte, 1), // Small buffer
	}

	// Fill the buffer
	conn.send <- []byte("first")

	// Try to send when buffer is full
	msg := &Message{
		ID:   uuid.New().String(),
		Type: "test",
	}

	err := conn.SendMessage(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "buffer full")
}
