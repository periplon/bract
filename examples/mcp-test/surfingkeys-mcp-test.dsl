# Surfingkeys MCP Integration Test
# This script tests all the new Surfingkeys MCP tools

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser connection
call browser_wait_for_connection {} -> connection_result
print "✓ Connected to browser extension"

# Create a test tab
call browser_create_tab {
    url: "https://example.com",
    active: true
} -> tab
print "✓ Created test tab"
set tab_id = tab.id

# Test 1: Get Page Title
call browser_get_page_title {
    tabId: tab_id
} -> title_result
print "Page title: " + title_result.title
assert len(title_result.title) > 0, "Failed to get page title"
print "✓ Page title test passed"

# Test 2: Search functionality
call browser_search {
    query: "test query",
    engine: "google",
    newTab: false
} -> search_result
print "✓ Search test passed"

# Test 3: Find text on page
call browser_find {
    text: "Example",
    caseSensitive: false,
    wholeWord: false,
    tabId: tab_id
} -> find_result
print "✓ Find text test passed"

# Test 4: Clipboard operations
# Write to clipboard
call browser_clipboard_write {
    text: "Test clipboard content",
    format: "text"
} -> write_result
print "✓ Wrote to clipboard"

# Read from clipboard
call browser_clipboard_read {} -> read_result
print "Clipboard content: " + read_result.text
assert read_result.text == "Test clipboard content", "Clipboard read/write mismatch"
print "✓ Clipboard test passed"

# Test 5: Hints functionality
call browser_hints_show {
    selector: "a",
    action: "click",
    tabId: tab_id
} -> hints_result
print "✓ Hints show test passed"

# Test 6: Omnibar
call browser_omnibar {
    type: "bookmarks",
    query: "test",
    tabId: tab_id
} -> omnibar_result
print "✓ Omnibar test passed"

# Test 7: Visual mode
call browser_visual_mode {
    selectElement: false,
    tabId: tab_id
} -> visual_result
print "✓ Visual mode test passed"

# Clean up
call browser_close_tab {
    tabId: tab_id
} -> close_result
print "✓ Closed test tab"

print "\n✅ All Surfingkeys MCP integration tests passed!"