# Multi-Tab Test
# Tests managing multiple browser tabs

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

# Create multiple tabs
print "Creating 3 tabs..."

# Create tabs individually (workaround for array append limitation)
call browser_create_tab -> tab1
print "Created tab 1: " + str(tab1.id)

call browser_create_tab -> tab2
print "Created tab 2: " + str(tab2.id)

call browser_create_tab -> tab3
print "Created tab 3: " + str(tab3.id)

# Store tabs in array
set tabs = [tab1, tab2, tab3]

# Navigate each tab to different pages
set urls = ["https://example.com", "https://example.org", "https://example.net"]

print "\nNavigating tabs..."
loop i in [0, 1, 2] {
  call browser_navigate {
    tabId: tabs[i].id,
    url: urls[i]
  }
  print "Tab " + str(i + 1) + " navigated to " + urls[i]
}

# List all tabs
call browser_list_tabs -> allTabs
print "\nTotal tabs open: " + str(len(allTabs))

# Activate middle tab
call browser_activate_tab {
  tabId: tabs[1].id
}
print "✓ Activated tab 2"

# Get title from each tab using browser_list_tabs
print "\nGetting titles from all tabs..."
call browser_list_tabs -> currentTabs
loop i in [0, 1, 2] {
  loop tab in currentTabs {
    if tab.id == tabs[i].id {
      print "Tab " + str(i + 1) + " title: " + tab.title
    }
  }
}

# Close tabs in reverse order
print "\nClosing tabs..."
loop i in [2, 1, 0] {
  call browser_close_tab {
    tabId: tabs[i].id
  }
  print "Closed tab " + str(i + 1)
}

print "\n✓ Multi-tab test completed successfully"