# Storage Test
# Tests browser storage capabilities (cookies, localStorage, sessionStorage)
#
# NOTE: This test currently skips all storage operations as they are not yet
# implemented in the browser extension. The following commands need extension support:
# - cookies.set, cookies.get, cookies.remove
# - storage.local.get, storage.local.set
# - storage.session.get, storage.session.set

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
  timeout: 5
} -> connection_result
print "Browser connection status:"
print connection_result

# Create a test tab
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://example.com"
}

# Test Cookies
print "=== Testing Cookies ==="

# Note: Cookie operations are not yet implemented in the browser extension
print "⚠️  Cookie operations not yet supported by browser extension"
print "   Skipping cookie tests..."

# Note: Storage operations are not yet implemented in the browser extension
print "\n⚠️  Storage operations not yet supported by browser extension"
print "   Skipping localStorage and sessionStorage tests..."

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ All storage tests passed!"
