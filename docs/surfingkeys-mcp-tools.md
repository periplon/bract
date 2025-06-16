# Surfingkeys MCP Integration Tools

This document describes the new MCP tools added to support Surfingkeys functionality.

## Overview

The Surfingkeys MCP integration adds support for advanced browser automation features including hints, search, clipboard operations, and visual mode. These tools enable programmatic control of Surfingkeys features through the MCP protocol.

## New Tools

### Hints Tools

#### browser_hints_show
Shows interactive element hints on the page.

**Parameters:**
- `selector` (string, optional): CSS selector to filter hints
- `action` (string, optional): Action to perform when hint is clicked (e.g., 'click', 'hover')
- `tabId` (number, optional): Tab ID (defaults to active tab)

**Example:**
```dsl
call browser_hints_show {
    selector: "a",
    action: "click",
    tabId: 123
}
```

#### browser_hints_click
Clicks on a hint element by selector, index, or text.

**Parameters:**
- `selector` (string, optional): CSS selector of the hint to click
- `index` (number, optional): Index of the hint to click (0-based)
- `text` (string, optional): Text content of the hint to click
- `tabId` (number, optional): Tab ID (defaults to active tab)

**Note:** At least one of `selector`, `index`, or `text` must be provided.

**Example:**
```dsl
call browser_hints_click {
    text: "Submit",
    tabId: 123
}
```

### Search Tools

#### browser_search
Performs a web search using the specified search engine.

**Parameters:**
- `query` (string, required): Search query
- `engine` (string, optional): Search engine to use (e.g., 'google', 'bing', 'duckduckgo')
- `newTab` (boolean, optional): Open search results in a new tab (default: true)

**Example:**
```dsl
call browser_search {
    query: "MCP protocol documentation",
    engine: "google",
    newTab: true
}
```

#### browser_find
Finds text on the current page.

**Parameters:**
- `text` (string, required): Text to find on the page
- `caseSensitive` (boolean, optional): Case sensitive search (default: false)
- `wholeWord` (boolean, optional): Match whole words only (default: false)
- `tabId` (number, optional): Tab ID (defaults to active tab)

**Example:**
```dsl
call browser_find {
    text: "MCP",
    caseSensitive: true,
    wholeWord: true,
    tabId: 123
}
```

### Clipboard Tools

#### browser_clipboard_read
Reads text from the system clipboard.

**Parameters:** None

**Returns:** JSON object with `text` field containing clipboard content

**Example:**
```dsl
call browser_clipboard_read {} -> result
print "Clipboard: " + result.text
```

#### browser_clipboard_write
Writes text to the system clipboard.

**Parameters:**
- `text` (string, required): Text to write to clipboard
- `format` (string, optional): Format of the clipboard content (default: 'text')

**Example:**
```dsl
call browser_clipboard_write {
    text: "Hello, clipboard!",
    format: "text"
}
```

### Other Tools

#### browser_omnibar
Shows the omnibar with specified type.

**Parameters:**
- `type` (string, required): Type of omnibar to show. Valid values: 'bookmarks', 'history', 'tabs', 'commands'
- `query` (string, optional): Initial query to populate in the omnibar
- `tabId` (number, optional): Tab ID (defaults to active tab)

**Example:**
```dsl
call browser_omnibar {
    type: "bookmarks",
    query: "mcp",
    tabId: 123
}
```

#### browser_visual_mode
Starts visual selection mode.

**Parameters:**
- `selectElement` (boolean, optional): Select entire element instead of text (default: false)
- `tabId` (number, optional): Tab ID (defaults to active tab)

**Example:**
```dsl
call browser_visual_mode {
    selectElement: true,
    tabId: 123
}
```

#### browser_get_page_title
Gets the title of the current page.

**Parameters:**
- `tabId` (number, optional): Tab ID (defaults to active tab)

**Returns:** JSON object with `title` field containing the page title

**Example:**
```dsl
call browser_get_page_title {
    tabId: 123
} -> result
print "Page title: " + result.title
```

## Integration with Existing Tools

These new tools work seamlessly with existing browser automation tools. For example:

1. Use `browser_create_tab` to open a new tab
2. Use `browser_search` to perform a search
3. Use `browser_hints_show` to display clickable hints
4. Use `browser_hints_click` to interact with elements
5. Use `browser_clipboard_write` to save results

## Testing

Test files are provided in `examples/mcp-test/`:
- `surfingkeys-mcp-test.dsl`: Comprehensive test of all new tools
- `surfingkeys-hints-test.dsl`: Detailed test of hints functionality

Run tests with:
```bash
dsl examples/mcp-test/surfingkeys-mcp-test.dsl
```

## Implementation Details

The implementation follows the existing pattern in the bract project:

1. **Browser Client Methods**: Added to `internal/browser/client.go`
2. **Handler Methods**: Added to `internal/handler/browser_handler.go`
3. **Tool Registration**: Added to `internal/mcp/server.go`
4. **Interface Updates**: Updated `internal/handler/interfaces.go`

All new methods follow the established patterns for error handling, parameter validation, and response formatting.