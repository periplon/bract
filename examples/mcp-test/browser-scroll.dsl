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
call browser_execute_script {
  tabId: tab.id,
  script: "JSON.stringify({x: window.pageXOffset, y: window.pageYOffset})"
} -> initialPos
print "Initial scroll position: " + initialPos

print "\n1. Testing scroll to specific coordinates:"
# Scroll to specific position
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 500,
  behavior: "instant"
} -> scrollResult1
print "✓ Scrolled to position (0, 500)"

# Verify scroll position
call browser_execute_script {
  tabId: tab.id,
  script: "window.pageYOffset"
} -> yPos1
assert yPos1 >= 400, "Should have scrolled down"
print "Current Y position: " + str(yPos1)

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

call browser_execute_script {
  tabId: tab.id,
  script: "window.pageYOffset"
} -> yPos2
print "Current Y position after smooth scroll: " + str(yPos2)

print "\n3. Testing scroll to element:"
# Find an element to scroll to (e.g., a heading)
call browser_scroll {
  tabId: tab.id,
  selector: "h2",
  behavior: "auto"
} -> scrollResult3
print "✓ Scrolled to first h2 element"

# Get the element's position to verify
call browser_execute_script {
  tabId: tab.id,
  script: "document.querySelector('h2').getBoundingClientRect().top"
} -> elementTop
print "Element distance from viewport top: " + str(elementTop)

print "\n4. Testing scroll without behavior (defaults to auto):"
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 0
} -> scrollResult4
print "✓ Scrolled back to top"

# Verify we're at the top
call browser_execute_script {
  tabId: tab.id,
  script: "window.pageYOffset"
} -> finalY
assert finalY <= 10, "Should be near the top"
print "Final Y position: " + str(finalY)

print "\n5. Testing horizontal scroll (if page allows):"
# Try horizontal scroll
call browser_scroll {
  tabId: tab.id,
  x: 100,
  y: 0,
  behavior: "instant"
} -> scrollResult5

call browser_execute_script {
  tabId: tab.id,
  script: "window.pageXOffset"
} -> xPos
print "Horizontal scroll position: " + str(xPos)

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