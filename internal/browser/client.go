package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/periplon/bract/internal/config"
)

// Client manages communication with the Chrome extension
type Client struct {
	config      config.WebSocketConfig
	connection  Connection
	mu          sync.RWMutex
	pending     sync.Map // map[string]chan Response
	activeTabID int
}

// Connection interface for WebSocket connection
type Connection interface {
	SendCommand(action string, data interface{}) (string, error)
}

// Response represents a response from the Chrome extension
type Response struct {
	Data  json.RawMessage
	Error string
}

// NewClient creates a new browser client
func NewClient(config config.WebSocketConfig) *Client {
	return &Client{
		config:      config,
		activeTabID: -1,
	}
}

// SetConnection sets the WebSocket connection
func (c *Client) SetConnection(conn Connection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connection = conn
}

// WaitForConnection waits for the WebSocket connection to be established
func (c *Client) WaitForConnection(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			c.mu.RLock()
			hasConnection := c.connection != nil
			c.mu.RUnlock()

			if hasConnection {
				return nil
			}

			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for Chrome extension connection")
			}
		}
	}
}

// RemoveConnection removes the WebSocket connection
func (c *Client) RemoveConnection(conn Connection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connection == conn {
		c.connection = nil
	}
}

// HandleResponse handles a response from the Chrome extension
func (c *Client) HandleResponse(id string, data json.RawMessage, errMsg string) {
	if ch, ok := c.pending.Load(id); ok {
		response := Response{
			Data:  data,
			Error: errMsg,
		}
		ch.(chan Response) <- response
		c.pending.Delete(id)
	}
}

// HandleEvent handles events from the Chrome extension
func (c *Client) HandleEvent(action string, data json.RawMessage) {
	// Handle browser events (e.g., tab closed, navigation)
	switch action {
	case "tabClosed":
		// Update active tab if necessary
		var event struct {
			TabID int `json:"tabId"`
		}
		if err := json.Unmarshal(data, &event); err == nil && event.TabID == c.activeTabID {
			c.activeTabID = -1
		}
	}
}

// sendCommand sends a command to the Chrome extension and waits for response
func (c *Client) sendCommand(ctx context.Context, action string, data interface{}) (json.RawMessage, error) {
	c.mu.RLock()
	conn := c.connection
	c.mu.RUnlock()

	if conn == nil {
		return nil, fmt.Errorf("no connection to Chrome extension")
	}

	// Send command
	msgID, err := conn.SendCommand(action, data)
	if err != nil {
		return nil, err
	}

	// Create response channel
	respChan := make(chan Response, 1)
	c.pending.Store(msgID, respChan)
	defer c.pending.Delete(msgID)

	// Wait for response with timeout
	select {
	case resp := <-respChan:
		if resp.Error != "" {
			return nil, fmt.Errorf("chrome extension error: %s", resp.Error)
		}
		return resp.Data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(time.Duration(c.config.ReconnectMs) * time.Millisecond):
		return nil, fmt.Errorf("timeout waiting for Chrome extension response")
	}
}

// Surfingkeys MCP Integration Methods

// ShowHints shows interactive element hints
func (c *Client) ShowHints(ctx context.Context, tabID int, selector, action string) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
	}

	if selector != "" {
		params["selector"] = selector
	}
	if action != "" {
		params["action"] = action
	}

	return c.sendCommand(ctx, "hints.show", params)
}

// ClickHint clicks on a hint element
func (c *Client) ClickHint(ctx context.Context, tabID int, selector string, index int, text string) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
	}

	if selector != "" {
		params["selector"] = selector
	}
	if index >= 0 {
		params["index"] = index
	}
	if text != "" {
		params["text"] = text
	}

	return c.sendCommand(ctx, "hints.click", params)
}

// Search performs a web search
func (c *Client) Search(ctx context.Context, query, engine string, newTab bool) (json.RawMessage, error) {
	params := map[string]interface{}{
		"query":  query,
		"newTab": newTab,
	}

	if engine != "" {
		params["engine"] = engine
	}

	return c.sendCommand(ctx, "search", params)
}

// Find searches for text on the current page
func (c *Client) Find(ctx context.Context, tabID int, text string, caseSensitive, wholeWord bool) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":         tabID,
		"text":          text,
		"caseSensitive": caseSensitive,
		"wholeWord":     wholeWord,
	}

	return c.sendCommand(ctx, "find", params)
}

// ReadClipboard reads the system clipboard
func (c *Client) ReadClipboard(ctx context.Context) (string, error) {
	data, err := c.sendCommand(ctx, "clipboard.read", nil)
	if err != nil {
		return "", err
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	return result.Text, nil
}

// WriteClipboard writes to the system clipboard
func (c *Client) WriteClipboard(ctx context.Context, text, format string) error {
	params := map[string]interface{}{
		"text": text,
	}

	if format != "" {
		params["format"] = format
	}

	_, err := c.sendCommand(ctx, "clipboard.write", params)
	return err
}

// ShowOmnibar shows the omnibar
func (c *Client) ShowOmnibar(ctx context.Context, tabID int, barType, query string) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
		"type":  barType,
	}

	if query != "" {
		params["query"] = query
	}

	return c.sendCommand(ctx, "omnibar.show", params)
}

// StartVisualMode starts visual selection mode
func (c *Client) StartVisualMode(ctx context.Context, tabID int, selectElement bool) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":         tabID,
		"selectElement": selectElement,
	}

	return c.sendCommand(ctx, "visual.start", params)
}

// GetPageTitle gets the title of the current page
func (c *Client) GetPageTitle(ctx context.Context, tabID int) (string, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
	}

	data, err := c.sendCommand(ctx, "getPageTitle", params)
	if err != nil {
		return "", err
	}

	var result struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	return result.Title, nil
}