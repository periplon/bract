# Surfingkeys Hints Test
# This script tests the hints functionality using only Surfingkeys tools

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser connection
call browser_wait_for_connection {} -> connection_result
print "✓ Connected to browser extension"

# Note: This test assumes you have a browser window open with a webpage loaded
# The test will use the current active tab

# Test 1: Show hints for all links
call browser_hints_show {
    selector: "a",
    action: "click"
} -> hints_all
print "✓ Showed hints for all links"

# Small delay
wait 2

# Test 2: Show hints for specific selector (if page has wiki links)
call browser_hints_show {
    selector: "a[href*='wiki']",
    action: "click"
} -> hints_wiki
print "✓ Showed hints for wiki links (if any)"

# Small delay
wait 2

# Test 3: Show hints with hover action
call browser_hints_show {
    selector: "button",
    action: "hover"
} -> hints_hover
print "✓ Showed hints with hover action (if buttons exist)"

# Test 4: Get page title
call browser_get_page_title {} -> title_result
print "Current page title: " + str(title_result)

# Test 5: Test clipboard operations
call browser_clipboard_write {
    text: "Test clipboard content from Surfingkeys"
} -> write_result
print "✓ Wrote to clipboard"

call browser_clipboard_read {} -> read_result
print "✓ Read from clipboard: " + str(read_result)

# Test 6: Test search functionality
call browser_search {
    query: "surfingkeys browser extension",
    engine: "google",
    newTab: true
} -> search_result
print "✓ Performed search (opened in new tab)"

# Wait a moment
wait 3

# Test 7: Find text on page
call browser_find {
    text: "browser",
    caseSensitive: false
} -> find_result
print "✓ Searched for text on page"

# Test 8: Show omnibar
call browser_omnibar {
    type: "bookmarks"
} -> omnibar_result
print "✓ Showed omnibar"

print "\n✅ All Surfingkeys tests completed!"
print "Note: Some tests require manual verification in the browser"