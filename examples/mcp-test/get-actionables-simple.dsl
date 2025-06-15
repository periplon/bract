# Simple Get Actionables Example
# Demonstrates the basic usage of browser_get_actionables tool

# Connect to the MCP browser server
connect "../../bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
  timeout: 30
} -> connection_result
print "Browser connection established"

# Create a tab and navigate to Wikipedia
call browser_create_tab {
  url: "https://en.wikipedia.org",
  active: true
} -> tab
print "Created tab with ID: " + str(tab.id)

# Wait for page to load
wait 2

# Get all actionable elements
call browser_get_actionables {
  tabId: tab.id
} -> actionables

# Display summary
print "\nTotal actionable elements: " + str(len(actionables))

# Show first 10 elements
print "\nFirst 10 actionable elements:"
set count = 0
loop item in actionables {
  if count < 10 {
    print str(count + 1) + ". [" + item.type + "] " + item.description
    print "   Selector: " + item.selector
    set count = count + 1
  }
}

# Count by type
set link_count = 0
set button_count = 0
set input_count = 0
set other_count = 0

loop item in actionables {
  if item.type == "a" {
    set link_count = link_count + 1
  }
  if item.type == "button" {
    set button_count = button_count + 1
  }
  if item.type == "input" || item.type == "textarea" || item.type == "select" || item.type == "input[type=\"checkbox\"]" || item.type == "input[type=\"radio\"]" {
    set input_count = input_count + 1
  }
  if item.type != "a" && item.type != "button" && item.type != "input" && item.type != "textarea" && item.type != "select" && item.type != "input[type=\"checkbox\"]" && item.type != "input[type=\"radio\"]" {
    set other_count = other_count + 1
  }
}

print "\nElements by type:"
print "- link: " + str(link_count)
print "- button: " + str(button_count)
print "- input: " + str(input_count)
print "- other: " + str(other_count)