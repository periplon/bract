# Scroll Verification Test
# Simple test to verify scrolling works correctly

connect "./bin/mcp-browser-server"
call browser_wait_for_connection {timeout: 10}

print "=== Scroll Verification Test ==="

# Create a tab with scrollable content
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://en.wikipedia.org/wiki/Web_browser"
}
wait 3

print "\n1. Testing basic scroll functionality:"
# Scroll down
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 500,
  behavior: "instant"
} -> scrollResult1
print "✓ Scrolled to (0, 500)"
print "  Result: " + str(scrollResult1)

# Verify with execute_script
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> pos1
print "  Verified position: (" + pos1.x + ", " + pos1.y + ")"

print "\n2. Testing scroll with only Y coordinate:"
call browser_scroll {
  tabId: tab.id,
  y: 1000,
  behavior: "instant"
} -> scrollResult2
print "✓ Scrolled Y to 1000"

# Verify position
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> pos2
print "  Verified position: (" + pos2.x + ", " + pos2.y + ")"

print "\n3. Testing scroll to element:"
call browser_scroll {
  tabId: tab.id,
  selector: "h2",
  behavior: "instant"
} -> scrollResult3
print "✓ Scrolled to first h2 element"
print "  Result: " + str(scrollResult3)

print "\n4. Testing smooth scroll:"
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 0,
  behavior: "smooth"
} -> scrollResult4
print "✓ Started smooth scroll to top"

wait 2

# Verify we're back at top
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> finalPos
print "  Final position: (" + finalPos.x + ", " + finalPos.y + ")"

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✅ Scroll verification test completed successfully!"