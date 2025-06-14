# Browser Navigation Test
# Tests browser automation capabilities

# Connect to the MCP browser server
connect "./mcp-browser-server"

# Create a new tab
call create_tab -> tab
print "Created tab:"
print tab

# Navigate to a website
call navigate {
  tabId: tab.id,
  url: "https://example.com"
} -> result

assert result.success == true, "Navigation failed"
print "✓ Successfully navigated to example.com"

# Wait for page to load
wait result.loaded == true, 5

# Get page title
call execute_script {
  tabId: tab.id,
  script: "document.title"
} -> title

print "Page title: " + title.result

# Take a screenshot
call screenshot {
  tabId: tab.id
} -> screenshot

assert screenshot.data != null, "Screenshot failed"
print "✓ Screenshot captured"

# Extract page content
call extract_content {
  tabId: tab.id,
  selector: "body"
} -> content

print "Page content length: " + str(len(content.text))

# Clean up - close the tab
call close_tab {
  tabId: tab.id
}

print "✓ Test completed successfully"