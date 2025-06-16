package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/periplon/bract/internal/browser"
)

// BrowserClient defines the interface for browser automation operations
type BrowserClient interface {
	// Connection management
	SetConnection(conn browser.Connection)
	RemoveConnection(conn browser.Connection)
	HandleResponse(id string, data json.RawMessage, errMsg string)
	HandleEvent(action string, data json.RawMessage)
	WaitForConnection(ctx context.Context, timeout time.Duration) error

	// Surfingkeys MCP Integration
	ShowHints(ctx context.Context, tabID int, selector, action string) (json.RawMessage, error)
	ClickHint(ctx context.Context, tabID int, selector string, index int, text string) (json.RawMessage, error)
	Search(ctx context.Context, query, engine string, newTab bool) (json.RawMessage, error)
	Find(ctx context.Context, tabID int, text string, caseSensitive, wholeWord bool) (json.RawMessage, error)
	ReadClipboard(ctx context.Context) (string, error)
	WriteClipboard(ctx context.Context, text, format string) error
	ShowOmnibar(ctx context.Context, tabID int, barType, query string) (json.RawMessage, error)
	StartVisualMode(ctx context.Context, tabID int, selectElement bool) (json.RawMessage, error)
	GetPageTitle(ctx context.Context, tabID int) (string, error)
}