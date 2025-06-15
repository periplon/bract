# Get Accessibility Snapshot of Current Tab
# Simple example showing how to analyze the accessibility tree of the active tab
#
# Prerequisites:
# 1. Browser with Perix extension installed
# 2. Extension connected to the WebSocket server
# 3. A web page loaded in at least one tab

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
print "Waiting for browser extension connection..."
call browser_wait_for_connection {timeout: 30} -> connection_result
print connection_result

print "\n=== Accessibility Snapshot of Current Tab ==="

# First check if there are any tabs
call browser_list_tabs -> tabs
set tabToUse = null

if len(tabs) == 0 {
  print "\nNo tabs found. Creating a new tab..."
  call browser_create_tab {
    url: "https://example.com",
    active: true
  } -> new_tab
  print "Created new tab with ID: " + str(new_tab.id)
  set tabToUse = new_tab
  
  # Wait for page to load
  print "Waiting for page to load..."
  call browser_wait_for_element {
    tabId: new_tab.id,
    selector: "body",
    timeout: 5000
  }
} else {
  # Use the first tab and activate it
  set tabToUse = tabs[0]
  print "\nFound " + str(len(tabs)) + " tab(s). Using tab ID: " + str(tabToUse.id)
  print "Tab URL: " + tabToUse.url
  
  # Activate the tab to ensure it's the current active tab
  call browser_activate_tab {
    tabId: tabToUse.id
  }
  print "Activated tab: " + str(tabToUse.id)
  
  # Wait a moment for the tab to be ready
  print "Ensuring page is loaded..."
  call browser_wait_for_element {
    tabId: tabToUse.id,
    selector: "body",
    timeout: 2000
  }
}

# Get accessibility snapshot of the specific tab
print "\nGetting accessibility snapshot for tab " + str(tabToUse.id) + "..."
call browser_get_accessibility_snapshot {
  tabId: tabToUse.id,
  interestingOnly: true
} -> snapshot_result

# The result is returned as a JSON string
print "\nResult received!"
set result_str = str(snapshot_result)
print "Length: " + str(len(result_str)) + " characters"

# Check if we got a valid snapshot
if result_str == "{\"snapshot\":null}" || result_str == "null" {
  print "\nWarning: Received null snapshot. This can happen when:"
  print "- The page is still loading"
  print "- The page has no accessible content"
  print "- The browser extension needs to be refreshed"
  print "\nTry refreshing the page or navigating to a content-rich website."
} else {
  print "\nAccessibility snapshot received successfully!"
  print "The snapshot contains the accessibility tree structure of the page."
  
  # Since we can't parse JSON in DSL, just show the length
  if len(result_str) > 500 {
    print "\nSnapshot size: " + str(len(result_str)) + " characters (large accessibility tree)"
  } else {
    print "\nSnapshot data:"
    print result_str
  }
}

print "\n✓ Example completed!"