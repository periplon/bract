# Browser Navigation Test
# Tests browser automation capabilities

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
  timeout: 5
} -> connection_result
print "Browser connection status:"
print connection_result

# Create a new tab
call browser_create_tab -> tab
print "Created tab:"
print tab

# Navigate to a website
call browser_navigate {
  tabId: tab.id,
  url: "https://example.com"
} -> result

assert result.success == true, "Navigation failed"
print "✓ Successfully navigated to example.com"

# Page is already loaded (navigate waits for load to complete)

# Get page title using browser_list_tabs
call browser_list_tabs -> tabs
set title = ""
loop t in tabs {
  if t.id == tab.id {
    set title = t.title
  }
}

print "Page title: " + title

# Take a screenshot
call browser_screenshot {
  tabId: tab.id
} -> screenshot

assert screenshot.dataUrl != null, "Screenshot failed"
print "✓ Screenshot captured"

# Extract page content
call browser_extract_content {
  tabId: tab.id,
  selector: "body"
} -> content

print "Page content length: " + str(len(content[0]))

# Clean up - close the tab
call browser_close_tab {
  tabId: tab.id
}

print "✓ Test completed successfully"
