# Accessibility Snapshot Example
# Demonstrates how to get the accessibility tree of a web page
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

# List existing tabs
call browser_list_tabs -> tabs

if len(tabs) == 0 {
  print "\nNo tabs found. Creating a new tab..."
  # Create a new tab if none exist
  call browser_create_tab {
    url: "https://www.w3.org/WAI/ARIA/apg/patterns/",
    active: true
  } -> tab
  print "Created tab: " + str(tab.id)
} else {
  # Use the first tab
  set tab = tabs[0]
  print "\nUsing existing tab: " + str(tab.id)
  
  # Navigate to the ARIA patterns page
  call browser_navigate {
    tabId: tab.id,
    url: "https://www.w3.org/WAI/ARIA/apg/patterns/"
  }
}

print "\n=== Accessibility Snapshot Example ==="

# Wait for page to load
print "Waiting for page to load..."
call browser_wait_for_element {
  tabId: tab.id,
  selector: "main",
  timeout: 5000
} -> wait_result

# Example 1: Get the full accessibility snapshot
print "\n1. Full Accessibility Snapshot (interesting nodes only):"
call browser_get_accessibility_snapshot {
  tabId: tab.id,
  interestingOnly: true
} -> snapshot

print "Successfully retrieved accessibility snapshot"
print "Snapshot length: " + str(len(str(snapshot))) + " characters"

# Example 2: Get a more detailed snapshot (all nodes)
print "\n2. Detailed Accessibility Snapshot (all nodes):"
call browser_get_accessibility_snapshot {
  tabId: tab.id,
  interestingOnly: false
} -> detailedSnapshot

print "Successfully retrieved detailed snapshot"
print "Detailed snapshot length: " + str(len(str(detailedSnapshot))) + " characters"

# Example 3: Get accessibility snapshot of a specific region
print "\n3. Accessibility Snapshot of Main Content:"
call browser_get_accessibility_snapshot {
  tabId: tab.id,
  interestingOnly: true,
  root: "main"
} -> mainSnapshot

print "Successfully retrieved main content snapshot"
print "Main content snapshot length: " + str(len(str(mainSnapshot))) + " characters"

# Show a small portion of the main content snapshot
print "\nMain content snapshot (truncated for display):"
print "Note: In a real application, you would parse this JSON data"

print "\n✓ Accessibility snapshot examples completed!"
print "\nThe accessibility snapshots contain JSON data with the following structure:"
print "- role: The ARIA role of the element"
print "- name: The accessible name/label" 
print "- children: Nested child elements"
print "- And other accessibility properties"