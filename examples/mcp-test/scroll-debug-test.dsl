# Debug Scroll Test
# Simple test to debug execute_script issues

connect "./bin/mcp-browser-server"
call browser_wait_for_connection {timeout: 10}

print "=== Debug Scroll Test ==="

# Create a simple test page
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "about:blank"
}
wait 1

print "\n1. Testing browser_execute_script return value:"
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: 100, y: 200 })"
} -> testResult

print "Raw result: " + str(testResult)

print "\n2. Testing window scroll values:"
call browser_execute_script {
  tabId: tab.id,
  script: "JSON.stringify({ x: window.scrollX, y: window.scrollY })"
} -> scrollPosString
print "Scroll position (string): " + scrollPosString

print "\n4. Testing scroll and return:"
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 0,
  behavior: "instant"
} -> scrollResult
print "Scroll result: " + str(scrollResult)
print "Scroll result type: " + type(scrollResult)

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ Debug test completed"