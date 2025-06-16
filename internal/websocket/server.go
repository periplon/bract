package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/periplon/bract/internal/browser"
)

// Server handles WebSocket connections from Chrome extensions
type Server struct {
	port           int
	browserClient  *browser.Client
	upgrader       websocket.Upgrader
	connections    sync.Map // map[string]*Connection
	allowedOrigins []string
}

// Connection represents a WebSocket connection to a Chrome extension
type Connection struct {
	ID        string
	conn      *websocket.Conn
	send      chan []byte
	server    *Server
	pingTimer *time.Timer
}

// Message represents a WebSocket message
type Message struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Command string          `json:"command,omitempty"` // Chrome extension expects 'command' field
	Action  string          `json:"action,omitempty"`  // Keep for backward compatibility
	Data    json.RawMessage `json:"data,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"` // Chrome extension sends responses with 'result' field
	Params  json.RawMessage `json:"params,omitempty"` // Chrome extension uses 'params' for data
	Error   string          `json:"error,omitempty"`
}

// NewServer creates a new WebSocket server
func NewServer(port int, browserClient *browser.Client, allowedOrigins []string) *Server {
	s := &Server{
		port:           port,
		browserClient:  browserClient,
		allowedOrigins: allowedOrigins,
	}

	s.upgrader = websocket.Upgrader{
		CheckOrigin: s.checkOrigin,
	}

	return s
}

// checkOrigin validates the origin of WebSocket connections
func (s *Server) checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Allow empty origin (same-origin requests)
	if origin == "" {
		return true
	}

	// Check against allowed origins from config
	for _, allowed := range s.allowedOrigins {
		// Handle wildcard Chrome extension origins
		if allowed == "chrome-extension://*" && strings.HasPrefix(origin, "chrome-extension://") {
			return true
		}

		// Handle localhost with any port
		if allowed == "http://localhost" && strings.HasPrefix(origin, "http://localhost") {
			return true
		}
		if allowed == "https://localhost" && strings.HasPrefix(origin, "https://localhost") {
			return true
		}

		// Exact match
		if origin == allowed {
			return true
		}
	}

	// Log rejected origins for debugging
	log.Printf("Rejected WebSocket connection from origin: %s", origin)
	return false
}

// Start starts the WebSocket server
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleWebSocket)
	mux.HandleFunc("/health", s.handleHealth)

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", s.port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	log.Printf("WebSocket server listening on ws://localhost:%d", s.port)
	return server.ListenAndServe()
}

// handleHealth provides a health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// handleWebSocket handles WebSocket upgrade and connection
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	connection := &Connection{
		ID:     uuid.New().String(),
		conn:   conn,
		send:   make(chan []byte, 256),
		server: s,
	}

	// Register connection
	s.connections.Store(connection.ID, connection)
	s.browserClient.SetConnection(connection)

	// Start connection handlers
	go connection.writePump()
	go connection.readPump()

	log.Printf("New WebSocket connection: %s", connection.ID)
}

// readPump handles incoming messages from the WebSocket connection
func (c *Connection) readPump() {
	defer func() {
		c.server.connections.Delete(c.ID)
		c.server.browserClient.RemoveConnection(c)
		c.conn.Close()
		close(c.send)
		if c.pingTimer != nil {
			c.pingTimer.Stop()
		}
		log.Printf("WebSocket connection closed: %s", c.ID)
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		log.Printf("Failed to set read deadline: %v", err)
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
			log.Printf("Failed to set read deadline in pong handler: %v", err)
		}
		return nil
	})

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Handle message based on type
		switch msg.Type {
		case "response":
			// Response to a command sent to the extension
			// Chrome extension sends response data in 'result' field
			responseData := msg.Result
			if responseData == nil {
				responseData = msg.Data // Fallback to 'data' field for backward compatibility
			}
			c.server.browserClient.HandleResponse(msg.ID, responseData, msg.Error)
		case "event":
			// Event from the extension (e.g., tab closed)
			c.server.browserClient.HandleEvent(msg.Action, msg.Data)
		case "ping":
			// Respond to ping
			if err := c.SendMessage(&Message{
				ID:      msg.ID,
				Type:    "pong",
				Command: "pong", // Add command field to satisfy Chrome extension validation
			}); err != nil {
				log.Printf("Failed to send pong message: %v", err)
			}
		case "connected":
			// Handle connection confirmation from Chrome extension
			log.Printf("Chrome extension connected successfully: %s", c.ID)
			// Only send acknowledgment if the incoming message has an ID
			if msg.ID != "" {
				if err := c.SendMessage(&Message{
					ID:      msg.ID,
					Type:    "ack",
					Command: "ack", // Add command field to satisfy Chrome extension validation
				}); err != nil {
					log.Printf("Failed to send acknowledgment: %v", err)
				}
			}
		case "error":
			// Handle error messages from Chrome extension
			log.Printf("Chrome extension error (connection: %s): %s", c.ID, msg.Error)
			if msg.Data != nil {
				log.Printf("Error details: %s", string(msg.Data))
			}
			// If this is a response to a command, handle it appropriately
			if msg.ID != "" {
				c.server.browserClient.HandleResponse(msg.ID, nil, msg.Error)
			}
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}

		// Reset read deadline
		if err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
			log.Printf("Failed to reset read deadline: %v", err)
		}
	}
}

// writePump handles outgoing messages to the WebSocket connection
func (c *Connection) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				log.Printf("Failed to set write deadline: %v", err)
				return
			}
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to write message: %v", err)
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
				log.Printf("Failed to set write deadline for ping: %v", err)
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage sends a message to the Chrome extension
func (c *Connection) SendMessage(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("connection send buffer full")
	}
}

// SendCommand sends a command to the Chrome extension and returns the message ID
func (c *Connection) SendCommand(action string, data interface{}) (string, error) {
	msgID := uuid.New().String()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Map action names to Chrome extension command format
	command := mapActionToCommand(action)

	msg := &Message{
		ID:      msgID,
		Type:    "command",
		Command: command,
		Action:  action,    // Keep for backward compatibility
		Data:    dataBytes, // Keep for backward compatibility
		Params:  dataBytes, // Chrome extension expects params
	}

	if err := c.SendMessage(msg); err != nil {
		return "", err
	}

	return msgID, nil
}

// mapActionToCommand maps Go action names to Chrome extension command names
func mapActionToCommand(action string) string {
	// For Surfingkeys integration, we pass the action as-is
	// The Surfingkeys commands are already in the correct format:
	// hints.show, hints.click, search, find, clipboard.read, clipboard.write,
	// omnibar.show, visual.start, getPageTitle
	return action
}

// GetConnections returns all active connections
func (s *Server) GetConnections() []*Connection {
	var connections []*Connection
	s.connections.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*Connection); ok {
			connections = append(connections, conn)
		}
		return true
	})
	return connections
}
