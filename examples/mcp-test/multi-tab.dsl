# Multi-Tab Test
# Tests managing multiple browser tabs

# Connect to the MCP browser server
connect "./mcp-browser-server"

# Create multiple tabs
set tabs = []

print "Creating 3 tabs..."
loop i in [1, 2, 3] {
  call create_tab -> tab
  set tabs = tabs + [tab]
  print "Created tab " + str(i) + ": " + tab.id
}

# Navigate each tab to different pages
set urls = [
  "https://example.com",
  "https://example.org",
  "https://example.net"
]

print "\nNavigating tabs..."
loop i in [0, 1, 2] {
  call navigate {
    tabId: tabs[i].id,
    url: urls[i]
  }
  print "Tab " + str(i + 1) + " navigated to " + urls[i]
}

# List all tabs
call list_tabs -> allTabs
print "\nTotal tabs open: " + str(len(allTabs))

# Activate middle tab
call activate_tab {
  tabId: tabs[1].id
}
print "✓ Activated tab 2"

# Get title from each tab
print "\nGetting titles from all tabs..."
loop i in [0, 1, 2] {
  call execute_script {
    tabId: tabs[i].id,
    script: "document.title"
  } -> title
  print "Tab " + str(i + 1) + " title: " + title
}

# Close tabs in reverse order
print "\nClosing tabs..."
loop i in [2, 1, 0] {
  call close_tab {
    tabId: tabs[i].id
  }
  print "Closed tab " + str(i + 1)
}

print "\n✓ Multi-tab test completed successfully"