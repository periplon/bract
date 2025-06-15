# Simple Accessibility Snapshot Test
# Tests the basic functionality of browser_get_accessibility_snapshot

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

print "=== Testing Accessibility Snapshot Tool ==="

# Get snapshot and just print it as-is to see what we get
call browser_get_accessibility_snapshot {} -> result

print "\nRaw result:"
print result

print "\n✓ Test completed!"