# Gmail Actionables Example
# Gets all actionable elements from Gmail and lists them categorized by type

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
  timeout: 30
} -> connection_result
print "Browser connection established"

# Create a tab and navigate to Gmail
call browser_create_tab {
  url: "https://mail.google.com",
  active: true
} -> tab
print "Created tab with ID: " + str(tab.id)

# Wait for page to load - Gmail may show login or inbox
call browser_wait_for_element {
  tabId: tab.id,
  selector: "input[type='email'], div[role='main']",
  timeout: 30000
} -> element_result
print "Gmail page loaded"

# Get all actionable elements
call browser_get_actionables {
  tabId: tab.id
} -> actionables
print "Total actionable elements found: " + str(len(actionables))

# Initialize categories
set compose_buttons = []
set email_items = []
set navigation_links = []
set search_inputs = []
set menu_buttons = []
set checkboxes = []
set other_buttons = []
set other_inputs = []
set other_links = []
set textareas = []
set selects = []

# Categorize actionables
loop item in actionables {
  # Check for compose/new message buttons
  if item.type == "button" || item.type == "a" {
    # Look for compose-related keywords in description
    if item.description == "Compose" || item.description == "New Message" || item.description == "New message" || item.description == "Compose new message" {
      set compose_buttons = compose_buttons + [item]
    }
  }
  
  # Categorize by type
  if item.type == "a" {
    set navigation_links = navigation_links + [item]
  }
  if item.type == "button" {
    set other_buttons = other_buttons + [item]
  }
  if item.type == "input" {
    # Check if it's a search input by common search-related descriptions
    if item.description == "Search" || item.description == "Search mail" || item.description == "Search in mail" {
      set search_inputs = search_inputs + [item]
    }
    set other_inputs = other_inputs + [item]
  }
  if item.type == "textarea" {
    set textareas = textareas + [item]
  }
  if item.type == "select" {
    set selects = selects + [item]
  }
  if item.type == "input[type=\"checkbox\"]" {
    set checkboxes = checkboxes + [item]
  }
}

# Display summary
print "\n=== Gmail Actionable Items Summary ==="
print "Total actionable elements: " + str(len(actionables))
print ""
print "Links: " + str(len(navigation_links))
print "Buttons: " + str(len(other_buttons))
print "Input fields: " + str(len(other_inputs))
print "Search inputs: " + str(len(search_inputs))
print "Textareas: " + str(len(textareas))
print "Select dropdowns: " + str(len(selects))
print "Checkboxes: " + str(len(checkboxes))

# Display key Gmail actions
print "\n=== Key Gmail Actions ==="

if len(compose_buttons) > 0 {
  print "\nCompose/New Message Buttons:"
  loop item in compose_buttons {
    print "  [" + str(item.labelNumber) + "] " + item.description
  }
}

if len(search_inputs) > 0 {
  print "\nSearch Inputs:"
  loop item in search_inputs {
    print "  [" + str(item.labelNumber) + "] " + item.description
  }
}

# Show first 10 buttons
if len(other_buttons) > 0 {
  print "\nButtons (first 10):"
  set count = 0
  loop item in other_buttons {
    if count < 10 {
      print "  [" + str(item.labelNumber) + "] " + item.description
      set count = count + 1
    }
  }
  if len(other_buttons) > 10 {
    print "  ... and " + str(len(other_buttons) - 10) + " more buttons"
  }
}

# Show first 10 links
if len(navigation_links) > 0 {
  print "\nLinks (first 10):"
  set count = 0
  loop item in navigation_links {
    if count < 10 {
      print "  [" + str(item.labelNumber) + "] " + item.description
      set count = count + 1
    }
  }
  if len(navigation_links) > 10 {
    print "  ... and " + str(len(navigation_links) - 10) + " more links"
  }
}

if len(checkboxes) > 0 {
  print "\nCheckboxes:"
  loop item in checkboxes {
    print "  [" + str(item.labelNumber) + "] " + item.description
  }
}

# Display detailed list of all elements
print "\n=== All Actionable Elements (Detailed) ==="
set displayed = 0
loop item in actionables {
  if displayed < 30 {  # Limit to first 30 for readability
    print "[" + str(item.labelNumber) + "] " + item.type + ": " + item.description
    if item.selector != "" {
      print "    Selector: " + item.selector
    }
    set displayed = displayed + 1
  }
}
if len(actionables) > 30 {
  print "\n... and " + str(len(actionables) - 30) + " more elements"
}

# Close the tab
call browser_close_tab {
  tabId: tab.id
} -> close_result

print "\n=== Analysis Complete ==="