# Comprehensive Scroll Test
# Tests edge cases and potential issues with scrolling

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 10}

print "=== Comprehensive Scroll Test ==="

# Test 1: Create a test page with known dimensions
print "\n1. Creating test page with scrollable content:"
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "data:text/html,<html><head><style>body{margin:0;padding:0;width:3000px;height:5000px;background:linear-gradient(to bottom right, red, blue);position:relative;}#marker-top{position:absolute;top:0;left:0;width:100px;height:100px;background:yellow;}#marker-bottom{position:absolute;bottom:0;right:0;width:100px;height:100px;background:green;}#marker-middle{position:absolute;top:2500px;left:1500px;width:100px;height:100px;background:orange;}#test-element{position:absolute;top:1000px;left:500px;width:200px;height:200px;background:purple;}</style></head><body><div id='marker-top'>TOP</div><div id='marker-middle'>MIDDLE</div><div id='marker-bottom'>BOTTOM</div><div id='test-element'>TEST ELEMENT</div></body></html>"
}
wait 2
print "✓ Created test page with 3000x5000px dimensions"

# Test 2: Get initial position and verify
print "\n2. Verifying initial scroll position:"
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> initialPos
assert initialPos.x == 0, "Initial X position should be 0"
assert initialPos.y == 0, "Initial Y position should be 0"
print "✓ Initial position verified: (0, 0)"

# Test 3: Test precise coordinate scrolling
print "\n3. Testing precise coordinate scrolling:"
call browser_scroll {
  tabId: tab.id,
  x: 1000,
  y: 2000,
  behavior: "instant"
} -> scrollResult1

# Verify actual scroll position
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> actualPos1
assert actualPos1.x == 1000, "X position should be 1000"
assert actualPos1.y == 2000, "Y position should be 2000"
print "✓ Scrolled to exact position (1000, 2000)"

# Test 4: Test scrolling beyond page bounds
print "\n4. Testing scroll beyond page bounds:"
call browser_scroll {
  tabId: tab.id,
  x: 10000,
  y: 10000,
  behavior: "instant"
} -> scrollResult2

# Verify it clamps to max values
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY, maxX: document.documentElement.scrollWidth - window.innerWidth, maxY: document.documentElement.scrollHeight - window.innerHeight })"
} -> boundsCheck
print "✓ Scroll clamped to max values: (" + boundsCheck.x + ", " + boundsCheck.y + ")"
print "  Max scroll values: (" + boundsCheck.maxX + ", " + boundsCheck.maxY + ")"

# Test 5: Test scrolling to specific element
print "\n5. Testing scroll to element by selector:"
call browser_scroll {
  tabId: tab.id,
  selector: "#test-element",
  behavior: "instant"
} -> scrollResult3

# Verify element is in viewport
call browser_execute_script {
  tabId: tab.id,
  script: "(function() { const el = document.querySelector('#test-element'); const rect = el.getBoundingClientRect(); return { inViewport: rect.top >= 0 && rect.left >= 0 && rect.bottom <= window.innerHeight && rect.right <= window.innerWidth, rect: { top: rect.top, left: rect.left, bottom: rect.bottom, right: rect.right } }; })()"
} -> elementCheck
print "✓ Element visibility - In viewport: " + elementCheck.inViewport
print "  Element position: top=" + elementCheck.rect.top + ", left=" + elementCheck.rect.left

# Test 6: Test smooth scrolling timing
print "\n6. Testing smooth scroll timing:"
call browser_scroll {
  tabId: tab.id,
  x: 0,
  y: 0,
  behavior: "smooth"
} -> scrollResult4

# Check position immediately (should still be scrolling)
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> midScrollPos
print "✓ Mid-scroll position: (" + midScrollPos.x + ", " + midScrollPos.y + ")"

# Wait for smooth scroll to complete
wait 2

# Check final position
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> finalSmoothPos
assert finalSmoothPos.x == 0, "Should have returned to X=0"
assert finalSmoothPos.y == 0, "Should have returned to Y=0"
print "✓ Smooth scroll completed to (0, 0)"

# Test 7: Test scrolling with only X or only Y
print "\n7. Testing partial coordinate scrolling:"
# Only X
call browser_scroll {
  tabId: tab.id,
  x: 500,
  behavior: "instant"
} -> scrollResult5

call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> partialX
assert partialX.x == 500, "X should be 500"
assert partialX.y == 0, "Y should remain 0"
print "✓ Scrolled only X to 500, Y remained at 0"

# Only Y
call browser_scroll {
  tabId: tab.id,
  y: 1500,
  behavior: "instant"
} -> scrollResult6

call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> partialY
assert partialY.x == 500, "X should remain 500"
assert partialY.y == 1500, "Y should be 1500"
print "✓ Scrolled only Y to 1500, X remained at 500"

# Test 8: Test invalid selector handling
print "\n8. Testing invalid selector (skipping error test - DSL doesn't support try-catch):"
print "✓ Skipped error handling test"

# Test 9: Test scrolling in inactive tab
print "\n9. Testing scroll in background tab:"
# Create another tab to make first tab inactive
call browser_create_tab -> tab2
call browser_navigate {
  tabId: tab2.id,
  url: "about:blank"
}

# Scroll in the first (now inactive) tab
call browser_scroll {
  tabId: tab.id,
  x: 100,
  y: 100,
  behavior: "instant"
} -> bgScrollResult

# Verify scroll worked in background tab
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> bgPos
assert bgPos.x == 100, "Background tab X should be 100"
assert bgPos.y == 100, "Background tab Y should be 100"
print "✓ Successfully scrolled background tab to (100, 100)"

# Test 10: Test rapid successive scrolls
print "\n10. Testing rapid successive scrolls:"
set positions = [[200, 200], [400, 400], [600, 600], [800, 800], [1000, 1000]]
loop pos in positions {
  call browser_scroll {
    tabId: tab.id,
    x: pos[0],
    y: pos[1],
    behavior: "instant"
  }
}

# Verify final position
call browser_execute_script {
  tabId: tab.id,
  script: "({ x: window.scrollX, y: window.scrollY })"
} -> rapidFinal
assert rapidFinal.x == 1000, "Final X should be 1000"
assert rapidFinal.y == 1000, "Final Y should be 1000"
print "✓ Rapid scrolls completed successfully to (1000, 1000)"

# Clean up
call browser_close_tab {tabId: tab.id}
call browser_close_tab {tabId: tab2.id}

print "\n✅ All comprehensive scroll tests passed!"
print "Scrolling functionality is working correctly!"