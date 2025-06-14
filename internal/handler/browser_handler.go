package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/periplon/bract/internal/browser"
	"time"
)

// BrowserHandler handles browser automation tool requests
type BrowserHandler struct {
	client BrowserClient
}

// NewBrowserHandler creates a new browser handler
func NewBrowserHandler(client BrowserClient) *BrowserHandler {
	return &BrowserHandler{
		client: client,
	}
}

// Connection Handlers

// WaitForConnection waits for the browser extension to connect
func (h *BrowserHandler) WaitForConnection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	timeout := time.Duration(request.GetFloat("timeout", 30)) * time.Second

	err := h.client.WaitForConnection(ctx, timeout)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to connect to browser: %v", err)), nil
	}

	return mcp.NewToolResultText("Successfully connected to browser extension"), nil
}

// Tab Management Handlers

// ListTabs lists all open browser tabs
func (h *BrowserHandler) ListTabs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tabs, err := h.client.ListTabs(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list tabs: %v", err)), nil
	}

	// Return tabs as JSON
	tabsJSON, err := json.Marshal(tabs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize tabs: %v", err)), nil
	}

	return mcp.NewToolResultText(string(tabsJSON)), nil
}

// CreateTab creates a new browser tab
func (h *BrowserHandler) CreateTab(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := request.GetString("url", "about:blank")
	active := request.GetBool("active", true)

	tab, err := h.client.CreateTab(ctx, url, active)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create tab: %v", err)), nil
	}

	// Return tab data as JSON for structured access
	tabJSON, err := json.Marshal(tab)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize tab data: %v", err)), nil
	}

	return mcp.NewToolResultText(string(tabJSON)), nil
}

// CloseTab closes a browser tab
func (h *BrowserHandler) CloseTab(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tabID, err := request.RequireInt("tabId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := h.client.CloseTab(ctx, tabID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to close tab: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Closed tab %d", tabID)), nil
}

// ActivateTab activates a browser tab
func (h *BrowserHandler) ActivateTab(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tabID, err := request.RequireInt("tabId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if err := h.client.ActivateTab(ctx, tabID); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to activate tab: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Activated tab %d", tabID)), nil
}

// Navigation Handlers

// Navigate navigates to a URL
func (h *BrowserHandler) Navigate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := request.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	waitUntilLoad := request.GetBool("waitUntilLoad", true)
	tabID := request.GetInt("tabId", 0)

	response, err := h.client.Navigate(ctx, tabID, url, waitUntilLoad)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to navigate: %v", err)), nil
	}

	// Return the browser extension's response (e.g., {success: true})
	return mcp.NewToolResultText(string(response)), nil
}

// Reload reloads the current page
func (h *BrowserHandler) Reload(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hardReload := request.GetBool("hardReload", false)
	tabID := request.GetInt("tabId", 0)

	if err := h.client.Reload(ctx, tabID, hardReload); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to reload: %v", err)), nil
	}

	reloadType := "Reloaded"
	if hardReload {
		reloadType = "Hard reloaded"
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s page", reloadType)), nil
}

// Interaction Handlers

// Click clicks on an element
func (h *BrowserHandler) Click(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, err := request.RequireString("selector")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	timeout := request.GetInt("timeout", 30000)
	tabID := request.GetInt("tabId", 0)

	if err := h.client.Click(ctx, tabID, selector, timeout); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to click: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Clicked on element: %s", selector)), nil
}

// Type types text into an input field
func (h *BrowserHandler) Type(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, err := request.RequireString("selector")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	text, err := request.RequireString("text")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clearFirst := request.GetBool("clearFirst", false)
	delay := request.GetInt("delay", 0)
	tabID := request.GetInt("tabId", 0)

	if err := h.client.Type(ctx, tabID, selector, text, clearFirst, delay); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to type: %v", err)), nil
	}

	action := "Typed"
	if clearFirst {
		action = "Cleared and typed"
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s '%s' into %s", action, text, selector)), nil
}

// Scroll scrolls the page
func (h *BrowserHandler) Scroll(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var x, y *float64
	if xVal := request.GetFloat("x", -1); xVal >= 0 {
		x = &xVal
	}
	if yVal := request.GetFloat("y", -1); yVal >= 0 {
		y = &yVal
	}

	selector := request.GetString("selector", "")
	behavior := request.GetString("behavior", "auto")
	tabID := request.GetInt("tabId", 0)

	response, err := h.client.Scroll(ctx, tabID, x, y, selector, behavior)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to scroll: %v", err)), nil
	}

	// If response contains scroll position data, return it
	if len(response) > 0 {
		return mcp.NewToolResultText(string(response)), nil
	}

	// Otherwise, return a descriptive message
	var scrollDesc string
	if selector != "" {
		scrollDesc = fmt.Sprintf("to element %s", selector)
	} else if x != nil && y != nil {
		scrollDesc = fmt.Sprintf("to position (%g, %g)", *x, *y)
	} else if x != nil {
		scrollDesc = fmt.Sprintf("horizontally to %g", *x)
	} else if y != nil {
		scrollDesc = fmt.Sprintf("vertically to %g", *y)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Scrolled %s", scrollDesc)), nil
}

// WaitForElement waits for an element to appear
func (h *BrowserHandler) WaitForElement(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector, err := request.RequireString("selector")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	timeout := request.GetInt("timeout", 30000)
	state := request.GetString("state", "visible")
	tabID := request.GetInt("tabId", 0)

	response, err := h.client.WaitForElement(ctx, tabID, selector, timeout, state)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to wait for element: %v", err)), nil
	}

	// If response contains element data, return it
	if len(response) > 0 {
		return mcp.NewToolResultText(string(response)), nil
	}

	// Otherwise, return a descriptive message
	return mcp.NewToolResultText(fmt.Sprintf("Element %s is now %s", selector, state)), nil
}

// Content Handlers

// ExecuteScript executes JavaScript in page context
func (h *BrowserHandler) ExecuteScript(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	script, err := request.RequireString("script")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get args as raw arguments and convert to slice
	var args []interface{}
	if rawArgs := request.GetArguments()["args"]; rawArgs != nil {
		if argsSlice, ok := rawArgs.([]interface{}); ok {
			args = argsSlice
		}
	}
	tabID := request.GetInt("tabId", 0)

	result, err := h.client.ExecuteScript(ctx, tabID, script, args)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to execute script: %v", err)), nil
	}

	// Return the raw script result as JSON
	return mcp.NewToolResultText(string(result)), nil
}

// ExtractContent extracts content from the page
func (h *BrowserHandler) ExtractContent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := request.GetString("selector", "body")
	contentType := request.GetString("type", "text")
	attribute := request.GetString("attribute", "")
	tabID := request.GetInt("tabId", 0)

	results, err := h.client.ExtractContent(ctx, tabID, selector, contentType, attribute)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to extract content: %v", err)), nil
	}

	// Return results as JSON array
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultsJSON)), nil
}

// Screenshot takes a screenshot
func (h *BrowserHandler) Screenshot(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fullPage := request.GetBool("fullPage", false)
	selector := request.GetString("selector", "")
	format := request.GetString("format", "png")
	quality := request.GetInt("quality", 90)
	tabID := request.GetInt("tabId", 0)

	dataURL, err := h.client.Screenshot(ctx, tabID, fullPage, selector, format, quality)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to take screenshot: %v", err)), nil
	}

	// Return as image content
	// Extract MIME type from data URL (format: data:image/png;base64,...)
	mimeType := "image/png" // default
	if len(dataURL) > 5 && strings.HasPrefix(dataURL, "data:") {
		if idx := strings.Index(dataURL, ";"); idx > 5 {
			mimeType = dataURL[5:idx]
		}
	}

	// Extract base64 data from data URL
	imageData := dataURL
	if idx := strings.Index(dataURL, ","); idx > 0 {
		imageData = dataURL[idx+1:]
	}

	return mcp.NewToolResultImage("Screenshot captured", imageData, mimeType), nil
}

// Storage Handlers

// GetCookies gets browser cookies
func (h *BrowserHandler) GetCookies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := request.GetString("url", "")
	name := request.GetString("name", "")

	cookies, err := h.client.GetCookies(ctx, url, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get cookies: %v", err)), nil
	}

	// Return cookies as JSON array
	cookiesJSON, err := json.Marshal(cookies)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize cookies: %v", err)), nil
	}

	return mcp.NewToolResultText(string(cookiesJSON)), nil
}

// SetCookie sets a browser cookie
func (h *BrowserHandler) SetCookie(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	value, err := request.RequireString("value")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	cookie := browser.Cookie{
		Name:           name,
		Value:          value,
		Domain:         request.GetString("domain", ""),
		Path:           request.GetString("path", "/"),
		Secure:         request.GetBool("secure", false),
		HTTPOnly:       request.GetBool("httpOnly", false),
		ExpirationDate: request.GetFloat("expirationDate", 0),
	}

	response, err := h.client.SetCookie(ctx, cookie)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set cookie: %v", err)), nil
	}

	// If response contains cookie data, return it
	if len(response) > 0 {
		return mcp.NewToolResultText(string(response)), nil
	}

	// Otherwise, return a descriptive message
	return mcp.NewToolResultText(fmt.Sprintf("Set cookie '%s' = '%s'", name, value)), nil
}

// DeleteCookies deletes browser cookies
func (h *BrowserHandler) DeleteCookies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := request.GetString("url", "")
	name := request.GetString("name", "")

	if err := h.client.DeleteCookies(ctx, url, name); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete cookies: %v", err)), nil
	}

	var desc string
	if name != "" {
		desc = fmt.Sprintf("cookie '%s'", name)
	} else if url != "" {
		desc = fmt.Sprintf("cookies for %s", url)
	} else {
		desc = "all cookies"
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted %s", desc)), nil
}

// GetLocalStorage gets localStorage value
func (h *BrowserHandler) GetLocalStorage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := request.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tabID := request.GetInt("tabId", 0)

	value, err := h.client.GetLocalStorage(ctx, tabID, key)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get localStorage: %v", err)), nil
	}

	// Return as JSON object
	result := map[string]string{"key": key, "value": value}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// SetLocalStorage sets localStorage value
func (h *BrowserHandler) SetLocalStorage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := request.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	value, err := request.RequireString("value")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tabID := request.GetInt("tabId", 0)

	if err := h.client.SetLocalStorage(ctx, tabID, key, value); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set localStorage: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set localStorage['%s'] = %s", key, value)), nil
}

// GetSessionStorage gets sessionStorage value
func (h *BrowserHandler) GetSessionStorage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := request.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tabID := request.GetInt("tabId", 0)

	value, err := h.client.GetSessionStorage(ctx, tabID, key)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get sessionStorage: %v", err)), nil
	}

	// Return as JSON object
	result := map[string]string{"key": key, "value": value}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// SetSessionStorage sets sessionStorage value
func (h *BrowserHandler) SetSessionStorage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := request.RequireString("key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	value, err := request.RequireString("value")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tabID := request.GetInt("tabId", 0)

	if err := h.client.SetSessionStorage(ctx, tabID, key, value); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to set sessionStorage: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Set sessionStorage['%s'] = %s", key, value)), nil
}
