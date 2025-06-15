# Simple Accessibility Snapshot Test
# Tests the browser_get_accessibility_snapshot tool with a known page

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
print "Waiting for browser extension..."
call browser_wait_for_connection {timeout: 30} -> connection_result
print connection_result

print "\n=== Testing Accessibility Snapshot Tool ==="

# Create a new tab with a simple page
print "\nCreating test tab..."
call browser_create_tab {
  url: "https://example.com",
  active: true
} -> tab
print "Created tab ID: " + str(tab.id)

# Wait for page to fully load
print "Waiting for page to load..."
call browser_wait_for_element {
  tabId: tab.id,
  selector: "h1",
  timeout: 5000
} -> wait_result
print "Page loaded"

# Test 1: Get accessibility snapshot with default settings
print "\n1. Testing with default settings (interesting nodes only)..."
call browser_get_accessibility_snapshot {
  tabId: tab.id
} -> snapshot1

set result1 = str(snapshot1)
print "Result length: " + str(len(result1)) + " characters"
if result1 == "{\"snapshot\":null}" || result1 == "null" {
  print "⚠️  Received null snapshot"
} else {
  print "✅ Received accessibility snapshot"
}

# Test 2: Get all nodes (not just interesting ones)
print "\n2. Testing with all nodes..."
call browser_get_accessibility_snapshot {
  tabId: tab.id,
  interestingOnly: false
} -> snapshot2

set result2 = str(snapshot2)
print "Result length: " + str(len(result2)) + " characters"
if result2 == "{\"snapshot\":null}" || result2 == "null" {
  print "⚠️  Received null snapshot"
} else {
  print "✅ Received accessibility snapshot"
}

# Test 3: Get snapshot of specific element
print "\n3. Testing with root selector (h1)..."
call browser_get_accessibility_snapshot {
  tabId: tab.id,
  root: "h1"
} -> snapshot3

set result3 = str(snapshot3)
print "Result length: " + str(len(result3)) + " characters"
if result3 == "{\"snapshot\":null}" || result3 == "null" {
  print "⚠️  Received null snapshot"
} else {
  print "✅ Received accessibility snapshot"
}

# Clean up
print "\nCleaning up..."
call browser_close_tab {
  tabId: tab.id
}
print "Closed test tab"

print "\n✓ Test completed!"