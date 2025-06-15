# Get Accessibility Snapshot of Current Tab
# Simple example showing how to analyze the accessibility tree of the active tab
#
# Prerequisites:
# 1. Browser with Perix extension installed
# 2. Extension connected to the WebSocket server

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
print "Waiting for browser extension connection..."
call browser_wait_for_connection {timeout: 30} -> connection_result
print connection_result

print "\n=== Accessibility Snapshot of Current Tab ==="

# Get accessibility snapshot of the current active tab (tabId defaults to 0 = active tab)
print "Getting accessibility snapshot..."
call browser_get_accessibility_snapshot {} -> snapshot_result

# The result is returned as a JSON string
print "\nResult received!"
print "Length: " + str(len(str(snapshot_result))) + " characters"

# Just print the full result - in a real app you'd parse this JSON
print "\nAccessibility snapshot:"
print snapshot_result

print "\n✓ Example completed!"
print "\nNote: If you see an error above, make sure:"
print "1. Browser with Perix extension is open"
print "2. Extension is connected (check extension popup)"
print "3. A web page is loaded in the active tab"