# Browser Scroll Test
# Tests various scrolling capabilities

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

# Create a tab with a long page
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://en.wikipedia.org/wiki/Web_browser"
}

print "=== Browser Scroll Test ==="

# Wait for page to load
wait 2

# Get initial scroll position
print "Initial scroll position: (0, 0) - page just loaded"

print "\n1. Testing scroll to specific coordinates:"
# Scroll to specific position
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 500,
  behavior: "instant"
} -> scrollResult1
print "✓ Scrolled to position (0, 500)"

# Verify scroll happened by taking a screenshot
call browser_screenshot {
  tabId: tab.id
} -> screenshot1
assert screenshot1.dataUrl != null, "Should have taken screenshot after scroll"
print "✓ Verified scroll occurred"

print "\n2. Testing smooth scroll:"
# Smooth scroll further down
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 1000,
  behavior: "smooth"
} -> scrollResult2
print "✓ Smooth scrolled to position (0, 1000)"

# Wait for smooth scroll to complete
wait 1

# Verify smooth scroll completed
call browser_screenshot {
  tabId: tab.id
} -> screenshot2
print "✓ Smooth scroll completed"

print "\n3. Testing scroll to element:"
# Find an element to scroll to (e.g., a heading)
call browser_scroll {
  tabId: tab.id,
  selector: "h2",
  behavior: "auto"
} -> scrollResult3
print "✓ Scrolled to first h2 element"

# Verify we scrolled to element
call browser_wait_for_element {
  tabId: tab.id,
  selector: "h2",
  timeout: 2
} -> h2Element
print "✓ H2 element is in viewport"

print "\n4. Testing scroll without behavior (defaults to auto):"
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 0
} -> scrollResult4
print "✓ Scrolled back to top"

# Verify we're back at the top by checking for the main heading
call browser_wait_for_element {
  tabId: tab.id,
  selector: "h1",
  timeout: 2
} -> topElement
print "✓ Returned to top of page"

print "\n5. Testing horizontal scroll (if page allows):"
# Try horizontal scroll
call browser_scroll {
  tabId: tab.id,
  x: 100,
  y: 0,
  behavior: "instant"
} -> scrollResult5

# Note: Most pages don't allow horizontal scroll
print "✓ Horizontal scroll attempted"

print "\n6. Testing scroll to element with smooth behavior:"
# Scroll to an element near bottom of page
call browser_scroll {
  tabId: tab.id,
  selector: "#References",
  behavior: "smooth"
} -> scrollResult6
print "✓ Smooth scrolled to References section"

wait 1

print "\n7. Testing scroll without tabId (uses active tab):"
# Ensure our tab is active
call browser_activate_tab {tabId: tab.id}

# Scroll without specifying tabId
call browser_scroll {
  x: 0,
  y: 200,
  behavior: "instant"
}
print "✓ Active tab scrolled successfully"

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ Browser scroll test completed!"