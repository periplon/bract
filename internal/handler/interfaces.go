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

	// Tab management
	ListTabs(ctx context.Context) ([]browser.Tab, error)
	CreateTab(ctx context.Context, url string, active bool) (*browser.Tab, error)
	CloseTab(ctx context.Context, tabID int) error
	ActivateTab(ctx context.Context, tabID int) error

	// Navigation
	Navigate(ctx context.Context, tabID int, url string, waitUntilLoad bool) (json.RawMessage, error)
	Reload(ctx context.Context, tabID int, hardReload bool) error

	// Interaction
	Click(ctx context.Context, tabID int, selector string, timeout int) error
	Type(ctx context.Context, tabID int, selector, text string, clearFirst bool, delay int) error
	Scroll(ctx context.Context, tabID int, x, y *float64, selector, behavior string) (json.RawMessage, error)
	WaitForElement(ctx context.Context, tabID int, selector string, timeout int, state string) (json.RawMessage, error)

	// Content
	ExecuteScript(ctx context.Context, tabID int, script string, args []interface{}) (json.RawMessage, error)
	ExtractContent(ctx context.Context, tabID int, selector, contentType, attribute string) ([]string, error)
	ExtractText(ctx context.Context, tabID int, selector string) (string, error)
	Screenshot(ctx context.Context, tabID int, fullPage bool, selector, format string, quality int) (string, error)

	// Storage
	GetCookies(ctx context.Context, url, name string) ([]browser.Cookie, error)
	SetCookie(ctx context.Context, cookie browser.Cookie) (json.RawMessage, error)
	DeleteCookies(ctx context.Context, url, name string) error
	GetLocalStorage(ctx context.Context, tabID int, key string) (string, error)
	SetLocalStorage(ctx context.Context, tabID int, key, value string) error
	ClearLocalStorage(ctx context.Context, tabID int) error
	GetSessionStorage(ctx context.Context, tabID int, key string) (string, error)
	SetSessionStorage(ctx context.Context, tabID int, key, value string) error
	ClearSessionStorage(ctx context.Context, tabID int) error

	// Actionables
	GetActionables(ctx context.Context, tabID int) ([]browser.Actionable, error)

	// Accessibility
	GetAccessibilitySnapshot(ctx context.Context, tabID int, interestingOnly bool, root string) (json.RawMessage, error)

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
