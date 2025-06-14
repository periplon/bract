# Basic Browser Connection Test
# This script tests the browser extension connection with wait

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
print "Waiting for browser extension to connect..."
call browser_wait_for_connection {
  timeout: 10
} -> connection_result
print "Connection result:"
print connection_result

# List available tools
call list_tools -> tools
print "Available browser tools:"
print tools

# Create a simple tab to verify connection works
call browser_create_tab {
  url: "about:blank",
  active: true
} -> tab
print "Successfully created tab:"
print tab

# Clean up
call browser_close_tab {
  tabId: tab.id
}
print "✓ Basic browser test completed"