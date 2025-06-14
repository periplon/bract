package handler

import (
	"context"
	"encoding/json"

	"github.com/periplon/bract/internal/browser"
)

// BrowserClient defines the interface for browser automation operations
type BrowserClient interface {
	// Connection management
	SetConnection(conn browser.Connection)
	RemoveConnection(conn browser.Connection)
	HandleResponse(id string, data json.RawMessage, errMsg string)
	HandleEvent(action string, data json.RawMessage)

	// Tab management
	ListTabs(ctx context.Context) ([]browser.Tab, error)
	CreateTab(ctx context.Context, url string, active bool) (*browser.Tab, error)
	CloseTab(ctx context.Context, tabID int) error
	ActivateTab(ctx context.Context, tabID int) error

	// Navigation
	Navigate(ctx context.Context, tabID int, url string, waitUntilLoad bool) error
	Reload(ctx context.Context, tabID int, hardReload bool) error

	// Interaction
	Click(ctx context.Context, tabID int, selector string, timeout int) error
	Type(ctx context.Context, tabID int, selector, text string, clearFirst bool, delay int) error
	Scroll(ctx context.Context, tabID int, x, y *float64, selector, behavior string) error
	WaitForElement(ctx context.Context, tabID int, selector string, timeout int, state string) error

	// Content
	ExecuteScript(ctx context.Context, tabID int, script string, args []interface{}) (json.RawMessage, error)
	ExtractContent(ctx context.Context, tabID int, selector, contentType, attribute string) ([]string, error)
	Screenshot(ctx context.Context, tabID int, fullPage bool, selector, format string, quality int) (string, error)

	// Storage
	GetCookies(ctx context.Context, url, name string) ([]browser.Cookie, error)
	SetCookie(ctx context.Context, cookie browser.Cookie) error
	DeleteCookies(ctx context.Context, url, name string) error
	GetLocalStorage(ctx context.Context, tabID int, key string) (string, error)
	SetLocalStorage(ctx context.Context, tabID int, key, value string) error
	GetSessionStorage(ctx context.Context, tabID int, key string) (string, error)
	SetSessionStorage(ctx context.Context, tabID int, key, value string) error
}