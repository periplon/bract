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
		var tabData struct {
			TabID int `json:"tabId"`
		}
		if err := json.Unmarshal(data, &tabData); err == nil {
			if tabData.TabID == c.activeTabID {
				c.activeTabID = -1
			}
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

// Tab Management Methods

// ListTabs returns all open browser tabs
func (c *Client) ListTabs(ctx context.Context) ([]Tab, error) {
	data, err := c.sendCommand(ctx, "listTabs", nil)
	if err != nil {
		return nil, err
	}

	var tabs []Tab
	if err := json.Unmarshal(data, &tabs); err != nil {
		return nil, err
	}

	return tabs, nil
}

// CreateTab creates a new browser tab
func (c *Client) CreateTab(ctx context.Context, url string, active bool) (*Tab, error) {
	params := map[string]interface{}{
		"url":    url,
		"active": active,
	}

	data, err := c.sendCommand(ctx, "createTab", params)
	if err != nil {
		return nil, err
	}

	// Check for empty response
	if len(data) == 0 {
		return nil, fmt.Errorf("empty response from browser extension - ensure Chrome extension is properly installed and running")
	}

	var tab Tab
	if err := json.Unmarshal(data, &tab); err != nil {
		return nil, fmt.Errorf("failed to parse tab response: %w (response: %s)", err, string(data))
	}

	if active {
		c.activeTabID = tab.ID
	}

	return &tab, nil
}

// CloseTab closes a browser tab
func (c *Client) CloseTab(ctx context.Context, tabID int) error {
	params := map[string]interface{}{
		"tabId": tabID,
	}

	_, err := c.sendCommand(ctx, "closeTab", params)
	return err
}

// ActivateTab activates (switches to) a browser tab
func (c *Client) ActivateTab(ctx context.Context, tabID int) error {
	params := map[string]interface{}{
		"tabId": tabID,
	}

	_, err := c.sendCommand(ctx, "activateTab", params)
	if err == nil {
		c.activeTabID = tabID
	}
	return err
}

// Navigation Methods

// Navigate navigates to a URL in a tab
func (c *Client) Navigate(ctx context.Context, tabID int, url string, waitUntilLoad bool) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":         tabID,
		"url":           url,
		"waitUntilLoad": waitUntilLoad,
	}

	response, err := c.sendCommand(ctx, "navigate", params)
	return response, err
}

// Reload reloads a tab
func (c *Client) Reload(ctx context.Context, tabID int, hardReload bool) error {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":      tabID,
		"hardReload": hardReload,
	}

	_, err := c.sendCommand(ctx, "reload", params)
	return err
}

// Interaction Methods

// Click clicks on an element
func (c *Client) Click(ctx context.Context, tabID int, selector string, timeout int) error {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":    tabID,
		"selector": selector,
		"timeout":  timeout,
	}

	_, err := c.sendCommand(ctx, "click", params)
	return err
}

// Type types text into an input field
func (c *Client) Type(ctx context.Context, tabID int, selector, text string, clearFirst bool, delay int) error {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":      tabID,
		"selector":   selector,
		"text":       text,
		"clearFirst": clearFirst,
		"delay":      delay,
	}

	_, err := c.sendCommand(ctx, "type", params)
	return err
}

// Scroll scrolls the page
func (c *Client) Scroll(ctx context.Context, tabID int, x, y *float64, selector, behavior string) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":    tabID,
		"behavior": behavior,
	}

	if x != nil {
		params["x"] = *x
	}
	if y != nil {
		params["y"] = *y
	}
	if selector != "" {
		params["selector"] = selector
	}

	response, err := c.sendCommand(ctx, "scroll", params)
	return response, err
}

// WaitForElement waits for an element to appear
func (c *Client) WaitForElement(ctx context.Context, tabID int, selector string, timeout int, state string) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":    tabID,
		"selector": selector,
		"timeout":  timeout,
		"state":    state,
	}

	response, err := c.sendCommand(ctx, "waitForElement", params)
	return response, err
}

// Content Methods

// ExecuteScript executes JavaScript in page context
func (c *Client) ExecuteScript(ctx context.Context, tabID int, script string, args []interface{}) (json.RawMessage, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":  tabID,
		"script": script,
		"args":   args,
	}

	return c.sendCommand(ctx, "executeScript", params)
}

// ExtractContent extracts content from the page
func (c *Client) ExtractContent(ctx context.Context, tabID int, selector, contentType, attribute string) ([]string, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":       tabID,
		"selector":    selector,
		"contentType": contentType,
	}

	if attribute != "" {
		params["attribute"] = attribute
	}

	data, err := c.sendCommand(ctx, "extractContent", params)
	if err != nil {
		return nil, err
	}

	var results []string
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Screenshot takes a screenshot
func (c *Client) Screenshot(ctx context.Context, tabID int, fullPage bool, selector, format string, quality int) (string, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":    tabID,
		"fullPage": fullPage,
		"format":   format,
		"quality":  quality,
	}

	if selector != "" {
		params["selector"] = selector
	}

	data, err := c.sendCommand(ctx, "screenshot", params)
	if err != nil {
		return "", err
	}

	var result struct {
		DataURL string `json:"dataUrl"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	return result.DataURL, nil
}

// Storage Methods

// GetCookies gets browser cookies
func (c *Client) GetCookies(ctx context.Context, url, name string) ([]Cookie, error) {
	params := map[string]interface{}{}
	if url != "" {
		params["url"] = url
	}
	if name != "" {
		params["name"] = name
	}

	data, err := c.sendCommand(ctx, "getCookies", params)
	if err != nil {
		return nil, err
	}

	var cookies []Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, err
	}

	return cookies, nil
}

// SetCookie sets a browser cookie
func (c *Client) SetCookie(ctx context.Context, cookie Cookie) (json.RawMessage, error) {
	response, err := c.sendCommand(ctx, "setCookie", cookie)
	return response, err
}

// DeleteCookies deletes browser cookies
func (c *Client) DeleteCookies(ctx context.Context, url, name string) error {
	params := map[string]interface{}{}
	if url != "" {
		params["url"] = url
	}
	if name != "" {
		params["name"] = name
	}

	_, err := c.sendCommand(ctx, "deleteCookies", params)
	return err
}

// GetLocalStorage gets localStorage value
func (c *Client) GetLocalStorage(ctx context.Context, tabID int, key string) (string, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
		"key":   key,
	}

	data, err := c.sendCommand(ctx, "getLocalStorage", params)
	if err != nil {
		return "", err
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return "", err
	}

	return value, nil
}

// SetLocalStorage sets localStorage value
func (c *Client) SetLocalStorage(ctx context.Context, tabID int, key, value string) error {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
		"key":   key,
		"value": value,
	}

	_, err := c.sendCommand(ctx, "setLocalStorage", params)
	return err
}

// GetSessionStorage gets sessionStorage value
func (c *Client) GetSessionStorage(ctx context.Context, tabID int, key string) (string, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
		"key":   key,
	}

	data, err := c.sendCommand(ctx, "getSessionStorage", params)
	if err != nil {
		return "", err
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return "", err
	}

	return value, nil
}

// SetSessionStorage sets sessionStorage value
func (c *Client) SetSessionStorage(ctx context.Context, tabID int, key, value string) error {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
		"key":   key,
		"value": value,
	}

	_, err := c.sendCommand(ctx, "setSessionStorage", params)
	return err
}
