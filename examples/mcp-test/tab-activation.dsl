# Tab Activation Test
# Tests switching between multiple browser tabs

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

print "=== Tab Activation Test ==="

# Create multiple tabs with different content
print "\nCreating 3 tabs with different pages..."

# Tab 1: Example.com
call browser_create_tab -> tab1
call browser_navigate {
  tabId: tab1.id,
  url: "https://example.com"
}
print "Created Tab 1: example.com (ID: " + str(tab1.id) + ")"

# Tab 2: Example.org
call browser_create_tab -> tab2
call browser_navigate {
  tabId: tab2.id,
  url: "https://example.org"
}
print "Created Tab 2: example.org (ID: " + str(tab2.id) + ")"

# Tab 3: Example.net
call browser_create_tab -> tab3
call browser_navigate {
  tabId: tab3.id,
  url: "https://example.net"
}
print "Created Tab 3: example.net (ID: " + str(tab3.id) + ")"

# List all tabs
call browser_list_tabs -> allTabs
print "\nTotal tabs open: " + str(len(allTabs))

print "\n1. Testing tab activation:"
# Activate tab 1
call browser_activate_tab {
  tabId: tab1.id
} -> activateResult1
print "✓ Activated Tab 1 (example.com)"

# Verify tab 1 is active by executing script
call browser_execute_script {
  tabId: tab1.id,
  script: "document.domain"
} -> domain1
assert domain1 == "example.com", "Tab 1 should be on example.com"
print "Confirmed active tab domain: " + domain1

# Wait a moment for user to see the tab switch
wait 1

print "\n2. Testing activation of middle tab:"
# Activate tab 2
call browser_activate_tab {
  tabId: tab2.id
} -> activateResult2
print "✓ Activated Tab 2 (example.org)"

# Perform an action on the active tab
call browser_execute_script {
  tabId: tab2.id,
  script: "document.title = 'Active Tab - ' + document.title; document.title"
} -> newTitle2
print "Modified active tab title: " + newTitle2

wait 1

print "\n3. Testing activation of last tab:"
# Activate tab 3
call browser_activate_tab {
  tabId: tab3.id
} -> activateResult3
print "✓ Activated Tab 3 (example.net)"

# Take a screenshot of the active tab
call browser_screenshot {
  tabId: tab3.id
} -> screenshot
assert screenshot.dataUrl != null, "Screenshot should be captured"
print "✓ Screenshot captured from active tab"

wait 1

print "\n4. Testing rapid tab switching:"
# Switch between tabs quickly
print "Rapidly switching between tabs..."

call browser_activate_tab {tabId: tab1.id}
wait 0.5
call browser_activate_tab {tabId: tab3.id}
wait 0.5
call browser_activate_tab {tabId: tab2.id}
wait 0.5
call browser_activate_tab {tabId: tab1.id}

print "✓ Rapid tab switching completed"

print "\n5. Testing operations on non-active tabs:"
# While tab1 is active, perform operations on tab2
print "Tab 1 is active, performing operations on Tab 2..."

call browser_execute_script {
  tabId: tab2.id,
  script: "localStorage.setItem('backgroundOp', 'success')"
}

call browser_get_local_storage {
  tabId: tab2.id,
  key: "backgroundOp"
} -> bgResult
assert bgResult.value == "success", "Background operation should succeed"
print "✓ Successfully performed operation on non-active tab"

print "\n6. Testing tab activation with page interaction:"
# Activate tab 3 and interact with it
call browser_activate_tab {tabId: tab3.id}

# Type something in the active tab's body (if there's an input)
call browser_execute_script {
  tabId: tab3.id,
  script: "document.body.innerHTML += '<p id=\"test\">Tab activated and modified!</p>'; 'done'"
}

# Verify the modification
call browser_extract_content {
  tabId: tab3.id,
  selector: "#test"
} -> testContent
print "✓ Active tab modified: " + testContent[0]

print "\n7. Verifying final state:"
# List tabs again to see current state
call browser_list_tabs -> finalTabs
print "Final tab count: " + str(len(finalTabs))

# Clean up - close all test tabs
print "\nCleaning up..."
call browser_close_tab {tabId: tab3.id}
print "Closed Tab 3"
call browser_close_tab {tabId: tab2.id}
print "Closed Tab 2"
call browser_close_tab {tabId: tab1.id}
print "Closed Tab 1"

print "\n✓ Tab activation test completed!"