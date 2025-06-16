# Debug Parameters Test
# Check what parameters are being sent

connect "./bin/mcp-browser-server"
call browser_wait_for_connection {timeout: 10}

print "=== Debug Parameters Test ==="

# Create a simple test page
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "about:blank"
}
wait 1

# Test 1: With all parameters
print "\n1. Testing with all parameters:"
call browser_scroll {
  tabId: tab.id,
  x: 100,
  y: 200,
  behavior: "instant"
}
print "✓ Scroll with all params completed"

# Test 2: With only y parameter
print "\n2. Testing with only y parameter:"
call browser_scroll {
  tabId: tab.id,
  y: 300
}
print "✓ Scroll with only y completed"

# Test 3: With selector
print "\n3. Testing with selector:"
call browser_navigate {
  tabId: tab.id,
  url: "data:text/html,<h1 id='test'>Test Header</h1><div style='height:2000px'></div>"
}
wait 1

call browser_scroll {
  tabId: tab.id,
  selector: "#test"
}
print "✓ Scroll with selector completed"

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ Debug test completed"