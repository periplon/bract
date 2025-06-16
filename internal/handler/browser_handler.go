package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
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

// Surfingkeys MCP Integration Handlers

// ShowHints shows interactive element hints
func (h *BrowserHandler) ShowHints(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := request.GetString("selector", "")
	action := request.GetString("action", "")
	tabID := request.GetInt("tabId", 0)

	response, err := h.client.ShowHints(ctx, tabID, selector, action)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to show hints: %v", err)), nil
	}

	return mcp.NewToolResultText(string(response)), nil
}

// ClickHint clicks on a hint element
func (h *BrowserHandler) ClickHint(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := request.GetString("selector", "")
	index := request.GetInt("index", -1)
	text := request.GetString("text", "")
	tabID := request.GetInt("tabId", 0)

	// At least one of selector, index, or text must be provided
	if selector == "" && index < 0 && text == "" {
		return mcp.NewToolResultError("Must provide either selector, index, or text parameter"), nil
	}

	response, err := h.client.ClickHint(ctx, tabID, selector, index, text)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to click hint: %v", err)), nil
	}

	return mcp.NewToolResultText(string(response)), nil
}

// Search performs a web search
func (h *BrowserHandler) Search(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	engine := request.GetString("engine", "")
	newTab := request.GetBool("newTab", true)

	response, err := h.client.Search(ctx, query, engine, newTab)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search: %v", err)), nil
	}

	return mcp.NewToolResultText(string(response)), nil
}

// Find searches for text on the current page
func (h *BrowserHandler) Find(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := request.RequireString("text")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	caseSensitive := request.GetBool("caseSensitive", false)
	wholeWord := request.GetBool("wholeWord", false)
	tabID := request.GetInt("tabId", 0)

	response, err := h.client.Find(ctx, tabID, text, caseSensitive, wholeWord)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to find text: %v", err)), nil
	}

	return mcp.NewToolResultText(string(response)), nil
}

// ReadClipboard reads the system clipboard
func (h *BrowserHandler) ReadClipboard(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := h.client.ReadClipboard(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read clipboard: %v", err)), nil
	}

	// Return as JSON object for consistency
	result := map[string]string{"text": text}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// WriteClipboard writes to the system clipboard
func (h *BrowserHandler) WriteClipboard(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := request.RequireString("text")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	format := request.GetString("format", "text")

	if err := h.client.WriteClipboard(ctx, text, format); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write clipboard: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Wrote to clipboard: %s", text)), nil
}

// ShowOmnibar shows the omnibar
func (h *BrowserHandler) ShowOmnibar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	barType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	query := request.GetString("query", "")
	tabID := request.GetInt("tabId", 0)

	response, err := h.client.ShowOmnibar(ctx, tabID, barType, query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to show omnibar: %v", err)), nil
	}

	return mcp.NewToolResultText(string(response)), nil
}

// StartVisualMode starts visual selection mode
func (h *BrowserHandler) StartVisualMode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selectElement := request.GetBool("selectElement", false)
	tabID := request.GetInt("tabId", 0)

	response, err := h.client.StartVisualMode(ctx, tabID, selectElement)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to start visual mode: %v", err)), nil
	}

	return mcp.NewToolResultText(string(response)), nil
}

// GetPageTitle gets the title of the current page
func (h *BrowserHandler) GetPageTitle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tabID := request.GetInt("tabId", 0)

	title, err := h.client.GetPageTitle(ctx, tabID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get page title: %v", err)), nil
	}

	// Return as JSON object for consistency
	result := map[string]string{"title": title}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}