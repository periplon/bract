package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/periplon/bract/internal/handler"
)

// Server wraps the MCP server with browser automation capabilities
type Server struct {
	mcpServer *server.MCPServer
	handler   *handler.BrowserHandler
}

// NewServer creates a new MCP server with browser automation tools
func NewServer(name, version string, h *handler.BrowserHandler) *Server {
	// Create MCP server with tool capabilities
	mcpServer := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	s := &Server{
		mcpServer: mcpServer,
		handler:   h,
	}

	// Register all browser automation tools
	s.registerTools()

	return s
}

// Start starts the MCP server using stdio transport
func (s *Server) Start() error {
	return server.ServeStdio(s.mcpServer)
}

// registerTools registers all browser automation tools
func (s *Server) registerTools() {
	// Connection Tools
	s.registerWaitForConnectionTool()

	// Tab Management Tools
	s.registerTabListTool()
	s.registerTabCreateTool()
	s.registerTabCloseTool()
	s.registerTabActivateTool()

	// Navigation Tools
	s.registerNavigateTool()
	s.registerReloadTool()

	// Interaction Tools
	s.registerClickTool()
	s.registerTypeTool()
	s.registerScrollTool()
	s.registerWaitForElementTool()

	// Content Tools
	s.registerExecuteScriptTool()
	s.registerExtractContentTool()
	s.registerExtractTextTool()
	s.registerScreenshotTool()
	s.registerGetActionablesTool()
	s.registerGetAccessibilitySnapshotTool()

	// Storage Tools
	s.registerCookieTools()
	s.registerStorageTools()

	// Surfingkeys MCP Integration Tools
	s.registerSurfingkeysTools()
}

// Connection Tools

func (s *Server) registerWaitForConnectionTool() {
	tool := mcp.NewTool("browser_wait_for_connection",
		mcp.WithDescription("Wait for the browser extension to connect"),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout in seconds (default: 30)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.WaitForConnection(ctx, request)
	})
}

// Tab Management Tools

func (s *Server) registerTabListTool() {
	tool := mcp.NewTool("browser_list_tabs",
		mcp.WithDescription("List all open browser tabs"),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ListTabs(ctx, request)
	})
}

func (s *Server) registerTabCreateTool() {
	tool := mcp.NewTool("browser_create_tab",
		mcp.WithDescription("Create a new browser tab"),
		mcp.WithString("url",
			mcp.Description("URL to open in the new tab"),
		),
		mcp.WithBoolean("active",
			mcp.Description("Whether to make the tab active"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.CreateTab(ctx, request)
	})
}

func (s *Server) registerTabCloseTool() {
	tool := mcp.NewTool("browser_close_tab",
		mcp.WithDescription("Close a browser tab"),
		mcp.WithNumber("tabId",
			mcp.Required(),
			mcp.Description("ID of the tab to close"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.CloseTab(ctx, request)
	})
}

func (s *Server) registerTabActivateTool() {
	tool := mcp.NewTool("browser_activate_tab",
		mcp.WithDescription("Activate (switch to) a browser tab"),
		mcp.WithNumber("tabId",
			mcp.Required(),
			mcp.Description("ID of the tab to activate"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ActivateTab(ctx, request)
	})
}

// Navigation Tools

func (s *Server) registerNavigateTool() {
	tool := mcp.NewTool("browser_navigate",
		mcp.WithDescription("Navigate to a URL in the active tab"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to navigate to"),
		),
		mcp.WithBoolean("waitUntilLoad",
			mcp.Description("Wait for page to fully load"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to navigate in (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.Navigate(ctx, request)
	})
}

func (s *Server) registerReloadTool() {
	tool := mcp.NewTool("browser_reload",
		mcp.WithDescription("Reload the current page"),
		mcp.WithBoolean("hardReload",
			mcp.Description("Force reload ignoring cache"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to reload (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.Reload(ctx, request)
	})
}

// Interaction Tools

func (s *Server) registerClickTool() {
	tool := mcp.NewTool("browser_click",
		mcp.WithDescription("Click on an element"),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector for the element to click"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout in milliseconds (default: 30000)"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to click in (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.Click(ctx, request)
	})
}

func (s *Server) registerTypeTool() {
	tool := mcp.NewTool("browser_type",
		mcp.WithDescription("Type text into an input field"),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector for the input field"),
		),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text to type"),
		),
		mcp.WithBoolean("clearFirst",
			mcp.Description("Clear the field before typing"),
		),
		mcp.WithNumber("delay",
			mcp.Description("Delay between keystrokes in ms"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to type in (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.Type(ctx, request)
	})
}

func (s *Server) registerScrollTool() {
	tool := mcp.NewTool("browser_scroll",
		mcp.WithDescription("Scroll the page"),
		mcp.WithNumber("x",
			mcp.Description("Horizontal scroll position"),
		),
		mcp.WithNumber("y",
			mcp.Description("Vertical scroll position"),
		),
		mcp.WithString("selector",
			mcp.Description("Element to scroll to"),
		),
		mcp.WithString("behavior",
			mcp.Description("Scroll behavior: auto, smooth, instant"),
			mcp.Enum("auto", "smooth", "instant"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to scroll in (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.Scroll(ctx, request)
	})
}

func (s *Server) registerWaitForElementTool() {
	tool := mcp.NewTool("browser_wait_for_element",
		mcp.WithDescription("Wait for an element to appear on the page"),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector for the element"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout in milliseconds (default: 30000)"),
		),
		mcp.WithString("state",
			mcp.Description("State to wait for: visible, hidden, attached, detached"),
			mcp.Enum("visible", "hidden", "attached", "detached"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to wait in (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.WaitForElement(ctx, request)
	})
}

// Content Tools

func (s *Server) registerExecuteScriptTool() {
	tool := mcp.NewTool("browser_execute_script",
		mcp.WithDescription("Execute JavaScript in page context"),
		mcp.WithString("script",
			mcp.Required(),
			mcp.Description("JavaScript code to execute"),
		),
		mcp.WithArray("args",
			mcp.Description("Arguments to pass to the script"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to execute in (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ExecuteScript(ctx, request)
	})
}

func (s *Server) registerExtractContentTool() {
	tool := mcp.NewTool("browser_extract_content",
		mcp.WithDescription("Extract content from the page"),
		mcp.WithString("selector",
			mcp.Description("CSS selector for element(s) to extract"),
		),
		mcp.WithString("type",
			mcp.Description("Type of content to extract: text, html, attribute"),
			mcp.Enum("text", "html", "attribute"),
		),
		mcp.WithString("attribute",
			mcp.Description("Attribute name (when type is 'attribute')"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to extract from (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ExtractContent(ctx, request)
	})
}

func (s *Server) registerExtractTextTool() {
	tool := mcp.NewTool("browser_extract_text",
		mcp.WithDescription("Extract content from the page and convert it to plain text"),
		mcp.WithString("selector",
			mcp.Description("CSS selector for element(s) to extract (defaults to 'body')"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to extract from (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ExtractText(ctx, request)
	})
}

func (s *Server) registerScreenshotTool() {
	tool := mcp.NewTool("browser_screenshot",
		mcp.WithDescription("Take a screenshot"),
		mcp.WithBoolean("fullPage",
			mcp.Description("Capture full page or just viewport"),
		),
		mcp.WithString("selector",
			mcp.Description("CSS selector for specific element"),
		),
		mcp.WithString("format",
			mcp.Description("Image format: png, jpeg"),
			mcp.Enum("png", "jpeg"),
		),
		mcp.WithNumber("quality",
			mcp.Description("JPEG quality (0-100)"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to capture (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.Screenshot(ctx, request)
	})
}

func (s *Server) registerGetActionablesTool() {
	tool := mcp.NewTool("browser_get_actionables",
		mcp.WithDescription("Get all actionable elements on the page (buttons, links, inputs, etc.)"),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to get actionables from (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.GetActionables(ctx, request)
	})
}

func (s *Server) registerGetAccessibilitySnapshotTool() {
	tool := mcp.NewTool("browser_get_accessibility_snapshot",
		mcp.WithDescription("Get the accessibility tree of the page for understanding page structure and elements"),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID to get snapshot from (defaults to active tab)"),
		),
		mcp.WithBoolean("interestingOnly",
			mcp.Description("Only return nodes with semantic meaning (default: true)"),
		),
		mcp.WithString("root",
			mcp.Description("CSS selector for the root element to start from (defaults to document body)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.GetAccessibilitySnapshot(ctx, request)
	})
}

// Storage Tools

func (s *Server) registerCookieTools() {
	// Get cookies
	getCookiesTool := mcp.NewTool("browser_get_cookies",
		mcp.WithDescription("Get browser cookies"),
		mcp.WithString("url",
			mcp.Description("URL to get cookies for"),
		),
		mcp.WithString("name",
			mcp.Description("Specific cookie name"),
		),
	)

	s.mcpServer.AddTool(getCookiesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.GetCookies(ctx, request)
	})

	// Set cookie
	setCookieTool := mcp.NewTool("browser_set_cookie",
		mcp.WithDescription("Set a browser cookie"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Cookie name"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("Cookie value"),
		),
		mcp.WithString("domain",
			mcp.Description("Cookie domain"),
		),
		mcp.WithString("path",
			mcp.Description("Cookie path"),
		),
		mcp.WithBoolean("secure",
			mcp.Description("Secure cookie flag"),
		),
		mcp.WithBoolean("httpOnly",
			mcp.Description("HttpOnly cookie flag"),
		),
		mcp.WithNumber("expirationDate",
			mcp.Description("Cookie expiration timestamp"),
		),
	)

	s.mcpServer.AddTool(setCookieTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.SetCookie(ctx, request)
	})

	// Delete cookies
	deleteCookiesTool := mcp.NewTool("browser_delete_cookies",
		mcp.WithDescription("Delete browser cookies"),
		mcp.WithString("url",
			mcp.Description("URL to delete cookies for"),
		),
		mcp.WithString("name",
			mcp.Description("Specific cookie name to delete"),
		),
	)

	s.mcpServer.AddTool(deleteCookiesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.DeleteCookies(ctx, request)
	})
}

func (s *Server) registerStorageTools() {
	// Local storage get
	getLocalStorageTool := mcp.NewTool("browser_get_local_storage",
		mcp.WithDescription("Get localStorage item"),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("Storage key"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(getLocalStorageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.GetLocalStorage(ctx, request)
	})

	// Local storage set
	setLocalStorageTool := mcp.NewTool("browser_set_local_storage",
		mcp.WithDescription("Set localStorage item"),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("Storage key"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("Storage value"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(setLocalStorageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.SetLocalStorage(ctx, request)
	})

	// Local storage clear
	clearLocalStorageTool := mcp.NewTool("browser_clear_local_storage",
		mcp.WithDescription("Clear all localStorage items"),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(clearLocalStorageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ClearLocalStorage(ctx, request)
	})

	// Session storage get
	getSessionStorageTool := mcp.NewTool("browser_get_session_storage",
		mcp.WithDescription("Get sessionStorage item"),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("Storage key"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(getSessionStorageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.GetSessionStorage(ctx, request)
	})

	// Session storage set
	setSessionStorageTool := mcp.NewTool("browser_set_session_storage",
		mcp.WithDescription("Set sessionStorage item"),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("Storage key"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("Storage value"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(setSessionStorageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.SetSessionStorage(ctx, request)
	})

	// Session storage clear
	clearSessionStorageTool := mcp.NewTool("browser_clear_session_storage",
		mcp.WithDescription("Clear all sessionStorage items"),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(clearSessionStorageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ClearSessionStorage(ctx, request)
	})
}

// Surfingkeys MCP Integration Tools

func (s *Server) registerSurfingkeysTools() {
	// Hints Tools
	s.registerShowHintsTool()
	s.registerClickHintTool()

	// Search Tools
	s.registerSearchTool()
	s.registerFindTool()

	// Clipboard Tools
	s.registerClipboardReadTool()
	s.registerClipboardWriteTool()

	// Other Tools
	s.registerOmnibarTool()
	s.registerVisualModeTool()
	s.registerGetPageTitleTool()
}

func (s *Server) registerShowHintsTool() {
	tool := mcp.NewTool("browser_hints_show",
		mcp.WithDescription("Show interactive element hints on the page"),
		mcp.WithString("selector",
			mcp.Description("CSS selector to filter hints (optional)"),
		),
		mcp.WithString("action",
			mcp.Description("Action to perform when hint is clicked (e.g., 'click', 'hover')"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ShowHints(ctx, request)
	})
}

func (s *Server) registerClickHintTool() {
	tool := mcp.NewTool("browser_hints_click",
		mcp.WithDescription("Click on a hint element by selector, index, or text"),
		mcp.WithString("selector",
			mcp.Description("CSS selector of the hint to click"),
		),
		mcp.WithNumber("index",
			mcp.Description("Index of the hint to click (0-based)"),
		),
		mcp.WithString("text",
			mcp.Description("Text content of the hint to click"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ClickHint(ctx, request)
	})
}

func (s *Server) registerSearchTool() {
	tool := mcp.NewTool("browser_search",
		mcp.WithDescription("Perform a web search using the specified search engine"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query"),
		),
		mcp.WithString("engine",
			mcp.Description("Search engine to use (e.g., 'google', 'bing', 'duckduckgo')"),
		),
		mcp.WithBoolean("newTab",
			mcp.Description("Open search results in a new tab (default: true)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.Search(ctx, request)
	})
}

func (s *Server) registerFindTool() {
	tool := mcp.NewTool("browser_find",
		mcp.WithDescription("Find text on the current page"),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text to find on the page"),
		),
		mcp.WithBoolean("caseSensitive",
			mcp.Description("Case sensitive search (default: false)"),
		),
		mcp.WithBoolean("wholeWord",
			mcp.Description("Match whole words only (default: false)"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.Find(ctx, request)
	})
}

func (s *Server) registerClipboardReadTool() {
	tool := mcp.NewTool("browser_clipboard_read",
		mcp.WithDescription("Read text from the system clipboard"),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ReadClipboard(ctx, request)
	})
}

func (s *Server) registerClipboardWriteTool() {
	tool := mcp.NewTool("browser_clipboard_write",
		mcp.WithDescription("Write text to the system clipboard"),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text to write to clipboard"),
		),
		mcp.WithString("format",
			mcp.Description("Format of the clipboard content (default: 'text')"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.WriteClipboard(ctx, request)
	})
}

func (s *Server) registerOmnibarTool() {
	tool := mcp.NewTool("browser_omnibar",
		mcp.WithDescription("Show the omnibar with specified type"),
		mcp.WithString("type",
			mcp.Required(),
			mcp.Description("Type of omnibar to show (e.g., 'bookmarks', 'history', 'tabs')"),
			mcp.Enum("bookmarks", "history", "tabs", "commands"),
		),
		mcp.WithString("query",
			mcp.Description("Initial query to populate in the omnibar"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.ShowOmnibar(ctx, request)
	})
}

func (s *Server) registerVisualModeTool() {
	tool := mcp.NewTool("browser_visual_mode",
		mcp.WithDescription("Start visual selection mode"),
		mcp.WithBoolean("selectElement",
			mcp.Description("Select entire element instead of text (default: false)"),
		),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.StartVisualMode(ctx, request)
	})
}

func (s *Server) registerGetPageTitleTool() {
	tool := mcp.NewTool("browser_get_page_title",
		mcp.WithDescription("Get the title of the current page"),
		mcp.WithNumber("tabId",
			mcp.Description("Tab ID (defaults to active tab)"),
		),
	)

	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handler.GetPageTitle(ctx, request)
	})
}
