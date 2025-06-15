# Gmail Actionables Example
# Gets all actionable elements from Gmail and lists them categorized by type

# Connect to the MCP browser server
connect "../../bin/mcp-browser-server"

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

# Categorize actionables
loop item in actionables {
  set desc_lower = lower(item.description)
  
  if item.type == "button" || item.type == "a" {
    if contains(desc_lower, "compose") || contains(desc_lower, "new message") {
      set compose_buttons = compose_buttons + [item]
    }
    if contains(desc_lower, "menu") || contains(desc_lower, "settings") || contains(desc_lower, "more") {
      set menu_buttons = menu_buttons + [item]
    }
    if item.type == "button" && !contains(desc_lower, "compose") && !contains(desc_lower, "new message") && !contains(desc_lower, "menu") && !contains(desc_lower, "settings") && !contains(desc_lower, "more") {
      set other_buttons = other_buttons + [item]
    }
    if item.type == "a" && !contains(desc_lower, "compose") && !contains(desc_lower, "new message") {
      set other_links = other_links + [item]
    }
  }
  if item.type == "input" {
    if contains(desc_lower, "search") {
      set search_inputs = search_inputs + [item]
    }
    if !contains(desc_lower, "search") {
      set other_inputs = other_inputs + [item]
    }
  }
  if item.type == "input[type=\"checkbox\"]" {
    set checkboxes = checkboxes + [item]
  }
  if contains(str(item.selector), "role=\"link\"") {
    set email_items = email_items + [item]
  }
  if item.type == "a" && !contains(str(item.selector), "role=\"link\"") {
    set navigation_links = navigation_links + [item]
  }
}

# Display summary
print "\n=== Gmail Actionable Items Summary ==="
print "Total actionable elements: " + str(len(actionables))
print ""
print "Compose Buttons: " + str(len(compose_buttons))
print "Email Items: " + str(len(email_items))
print "Navigation Links: " + str(len(navigation_links))
print "Search Inputs: " + str(len(search_inputs))
print "Menu Buttons: " + str(len(menu_buttons))
print "Checkboxes: " + str(len(checkboxes))
print "Other Buttons: " + str(len(other_buttons))
print "Other Inputs: " + str(len(other_inputs))
print "Other Links: " + str(len(other_links))

# Display key Gmail actions
print "\n=== Key Gmail Actions ==="

if len(compose_buttons) > 0 {
  print "\nCompose/New Message:"
  loop item in compose_buttons {
    print "  [" + str(item.labelNumber) + "] " + item.description
  }
}

if len(search_inputs) > 0 {
  print "\nSearch:"
  loop item in search_inputs {
    print "  [" + str(item.labelNumber) + "] " + item.description
  }
}

if len(email_items) > 0 {
  print "\nEmail Items (clickable emails):"
  # Show first 5 email items
  set count = 0
  loop item in email_items {
    if count < 5 {
      print "  [" + str(item.labelNumber) + "] " + item.description
      set count = count + 1
    }
  }
  if len(email_items) > 5 {
    print "  ... and " + str(len(email_items) - 5) + " more emails"
  }
}

if len(checkboxes) > 0 {
  print "\nSelection Checkboxes:"
  loop item in checkboxes {
    print "  [" + str(item.labelNumber) + "] " + item.description
  }
}

# Display all actionables with details
print "\n=== All Actionable Elements ==="
loop item in actionables {
  print "[" + str(item.labelNumber) + "] " + item.type + ": " + item.description
  if item.selector != "" {
    print "    Selector: " + item.selector
  }
}

# Close the tab
call browser_close_tab {
  tabId: tab.id
} -> close_result

print "\n=== Analysis Complete ==="