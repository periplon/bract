# Surfingkeys Hints Test
# This script tests the hints functionality in detail

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser connection
call browser_wait_for_connection {} -> connection_result
print "✓ Connected to browser extension"

# Create a test tab with a page that has many links
call browser_create_tab {
    url: "https://www.wikipedia.org",
    active: true
} -> tab
print "✓ Created test tab"
set tab_id = tab.id

# Wait for page to load
call browser_wait_for_element {
    selector: "body",
    timeout: 10000,
    state: "visible",
    tabId: tab_id
} -> wait_result
print "✓ Page loaded"

# Test 1: Show hints for all links
call browser_hints_show {
    selector: "a",
    action: "click",
    tabId: tab_id
} -> hints_all
print "✓ Showed hints for all links"

# Test 2: Show hints for specific selector
call browser_hints_show {
    selector: "a[href*='wiki']",
    action: "click",
    tabId: tab_id
} -> hints_wiki
print "✓ Showed hints for wiki links"

# Test 3: Click hint by index
call browser_hints_click {
    index: 0,
    tabId: tab_id
} -> click_index
print "✓ Clicked hint by index"

# Navigate back
call browser_execute_script {
    script: "window.history.back()",
    tabId: tab_id
} -> back_result

# Wait for page to reload
wait 2

# Test 4: Click hint by text
call browser_hints_click {
    text: "English",
    tabId: tab_id
} -> click_text
print "✓ Clicked hint by text"

# Navigate back again
call browser_execute_script {
    script: "window.history.back()",
    tabId: tab_id
} -> back_result2

wait 2

# Test 5: Show hints with hover action
call browser_hints_show {
    selector: "button",
    action: "hover",
    tabId: tab_id
} -> hints_hover
print "✓ Showed hints with hover action"

# Clean up
call browser_close_tab {
    tabId: tab_id
} -> close_result
print "✓ Closed test tab"

print "\n✅ All hints tests passed!"