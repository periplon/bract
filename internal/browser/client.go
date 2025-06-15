package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
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

// SendKey sends keyboard key events to the browser
func (c *Client) SendKey(ctx context.Context, key string, modifiers map[string]bool, tabID int) error {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":     tabID,
		"key":       key,
		"modifiers": modifiers,
	}

	_, err := c.sendCommand(ctx, "sendKey", params)
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

	// Browser extension returns {text: ...}
	var response struct {
		Text interface{} `json:"text"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	// Handle both string and array responses
	switch v := response.Text.(type) {
	case string:
		return []string{v}, nil
	case []interface{}:
		results := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				results[i] = str
			}
		}
		return results, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", response.Text)
	}
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

	// Browser extension returns {cookies: [...]}
	var response struct {
		Cookies []Cookie `json:"cookies"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Cookies, nil
}

// SetCookie sets a browser cookie
func (c *Client) SetCookie(ctx context.Context, cookie Cookie) (json.RawMessage, error) {
	// Chrome extension requires a URL parameter
	params := map[string]interface{}{
		"name":     cookie.Name,
		"value":    cookie.Value,
		"domain":   cookie.Domain,
		"path":     cookie.Path,
		"secure":   cookie.Secure,
		"httpOnly": cookie.HTTPOnly,
	}

	// Build URL from domain
	protocol := "http"
	if cookie.Secure {
		protocol = "https"
	}
	domain := cookie.Domain
	if domain == "" {
		domain = "localhost"
	}
	// Remove leading dot from domain if present
	if len(domain) > 0 && domain[0] == '.' {
		domain = domain[1:]
	}
	params["url"] = fmt.Sprintf("%s://%s", protocol, domain)

	if cookie.ExpirationDate > 0 {
		params["expirationDate"] = cookie.ExpirationDate
	}

	response, err := c.sendCommand(ctx, "setCookie", params)
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

	// Browser extension returns {storage: {key: value}}
	var response struct {
		Storage map[string]interface{} `json:"storage"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return "", err
	}

	// Get the specific key value
	if response.Storage != nil {
		if val, ok := response.Storage[key]; ok {
			if strVal, ok := val.(string); ok {
				return strVal, nil
			}
			// Handle null values
			if val == nil {
				return "", nil
			}
		}
	}

	return "", nil
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

// ClearLocalStorage clears all localStorage
func (c *Client) ClearLocalStorage(ctx context.Context, tabID int) error {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
	}

	_, err := c.sendCommand(ctx, "clearLocalStorage", params)
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

	// Browser extension returns {storage: {key: value}}
	var response struct {
		Storage map[string]interface{} `json:"storage"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return "", err
	}

	// Get the specific key value
	if response.Storage != nil {
		if val, ok := response.Storage[key]; ok {
			if strVal, ok := val.(string); ok {
				return strVal, nil
			}
			// Handle null values
			if val == nil {
				return "", nil
			}
		}
	}

	return "", nil
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

// ClearSessionStorage clears all sessionStorage
func (c *Client) ClearSessionStorage(ctx context.Context, tabID int) error {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
	}

	_, err := c.sendCommand(ctx, "clearSessionStorage", params)
	return err
}

// GetActionables gets all actionable elements on the page
func (c *Client) GetActionables(ctx context.Context, tabID int) ([]Actionable, error) {
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId": tabID,
	}

	data, err := c.sendCommand(ctx, "tabs.getActionables", params)
	if err != nil {
		return nil, err
	}

	// Chrome extension returns { actionables: [...] }
	var response struct {
		Actionables []Actionable `json:"actionables"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Actionables, nil
}

// ExtractText extracts content from the page as HTML and converts it to plain text
func (c *Client) ExtractText(ctx context.Context, tabID int, selector string) (string, error) {
	// Use the existing extractContent command with type "html"
	results, err := c.ExtractContent(ctx, tabID, selector, "html", "")
	if err != nil {
		return "", err
	}

	// Join all results and strip HTML tags
	var plainText strings.Builder
	for i, html := range results {
		if i > 0 {
			plainText.WriteString("\n\n")
		}
		plainText.WriteString(stripHTMLTags(html))
	}

	return plainText.String(), nil
}

// stripHTMLTags removes HTML tags from a string and returns plain text
func stripHTMLTags(html string) string {
	// Remove script and style tags and their contents
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`)
	html = scriptRegex.ReplaceAllString(html, "")
	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>[\s\S]*?</style>`)
	html = styleRegex.ReplaceAllString(html, "")

	// Replace br tags with newlines
	brRegex := regexp.MustCompile(`(?i)<br\s*/?>`)
	html = brRegex.ReplaceAllString(html, "\n")

	// Replace p, div, and other block tags with newlines
	blockRegex := regexp.MustCompile(`(?i)</?(p|div|h[1-6]|ul|ol|li|blockquote|pre|table|tr|td|th)[^>]*>`)
	html = blockRegex.ReplaceAllString(html, "\n")

	// Remove all remaining HTML tags
	tagRegex := regexp.MustCompile(`<[^>]+>`)
	text := tagRegex.ReplaceAllString(html, "")

	// Decode HTML entities
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&nbsp;", " ")

	// Clean up excessive whitespace
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

// GetAccessibilitySnapshot gets the accessibility tree of the page
func (c *Client) GetAccessibilitySnapshot(ctx context.Context, tabID int, interestingOnly bool, root string) (json.RawMessage, error) {
	// Default to active tab if not specified
	if tabID == 0 {
		tabID = c.activeTabID
	}

	params := map[string]interface{}{
		"tabId":           tabID,
		"interestingOnly": interestingOnly,
	}

	if root != "" {
		params["root"] = root
	}

	result, err := c.sendCommand(ctx, "tabs.getAccessibilitySnapshot", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
