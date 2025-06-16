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