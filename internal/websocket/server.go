package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/periplon/bract/internal/browser"
)

// Server handles WebSocket connections from Chrome extensions
type Server struct {
	port          int
	browserClient *browser.Client
	upgrader      websocket.Upgrader
	connections   sync.Map // map[string]*Connection
}

// Connection represents a WebSocket connection to a Chrome extension
type Connection struct {
	ID        string
	conn      *websocket.Conn
	send      chan []byte
	server    *Server
	mu        sync.Mutex
	pingTimer *time.Timer
}

// Message represents a WebSocket message
type Message struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Action string          `json:"action,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// NewServer creates a new WebSocket server
func NewServer(port int, browserClient *browser.Client) *Server {
	return &Server{
		port:          port,
		browserClient: browserClient,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Only allow connections from Chrome extensions
				origin := r.Header.Get("Origin")
				return origin == "" || 
					   origin == "chrome-extension://"+r.Header.Get("X-Extension-Id") ||
					   origin == "http://localhost" ||
					   origin == "https://localhost"
			},
		},
	}
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
	w.Write([]byte("OK"))
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

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
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
			c.server.browserClient.HandleResponse(msg.ID, msg.Data, msg.Error)
		case "event":
			// Event from the extension (e.g., tab closed)
			c.server.browserClient.HandleEvent(msg.Action, msg.Data)
		case "ping":
			// Respond to ping
			c.SendMessage(&Message{
				ID:   msg.ID,
				Type: "pong",
			})
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}

		// Reset read deadline
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
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
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
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

	msg := &Message{
		ID:     msgID,
		Type:   "command",
		Action: action,
		Data:   dataBytes,
	}

	if err := c.SendMessage(msg); err != nil {
		return "", err
	}

	return msgID, nil
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