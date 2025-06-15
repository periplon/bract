# Wikipedia Actionables Example
# This example demonstrates how to use the browser_get_actionables tool
# to discover all interactive elements on a Wikipedia page

# Connect to the MCP browser server
connect "../../bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
  timeout: 30
} -> connection_result
print "Browser connected"

# Create a new tab and navigate to Wikipedia
call browser_create_tab {
  url: "https://en.wikipedia.org/wiki/Main_Page",
  active: true
} -> tab
print "Created tab with ID: " + str(tab.id)

# Wait for the page to fully load
call browser_wait_for_element {
  tabId: tab.id,
  selector: "body",
  timeout: 10000
} -> wait_result

# Get all actionable elements on the Wikipedia main page
call browser_get_actionables {
  tabId: tab.id
} -> actionables

# Print the total number of actionable elements found
print "\nFound " + str(len(actionables)) + " actionable elements on Wikipedia's main page"
print "---"

# Display the first 20 actionable elements with their details
print "\nFirst 20 actionable elements:"
set counter = 0
loop item in actionables {
  if counter < 20 {
    print "\nElement #" + str(item.labelNumber) + ":"
    print "  Description: " + item.description
    print "  Type: " + item.type
    print "  Selector: " + item.selector
    set counter = counter + 1
  }
}

# Find and display specific types of elements
print "\n=== Links ==="
set linkCount = 0
loop item in actionables {
  if item.type == "a" && linkCount < 10 {
    print "- " + item.description + " (" + item.selector + ")"
    set linkCount = linkCount + 1
  }
}

print "\n=== Buttons ==="
set buttonCount = 0
loop item in actionables {
  if item.type == "button" && buttonCount < 10 {
    print "- " + item.description + " (" + item.selector + ")"
    set buttonCount = buttonCount + 1
  }
}

print "\n=== Input Fields ==="
set inputCount = 0
loop item in actionables {
  if item.type == "input" && inputCount < 10 {
    print "- " + item.description + " (" + item.selector + ")"
    set inputCount = inputCount + 1
  }
}

# Example: Find and interact with search functionality
print "\n=== Search Interaction Example ==="
set searchInput = null
set searchButton = null

# Find search input
loop item in actionables {
  if item.type == "input" {
    # Check if description or selector contains "search"
    set isSearch = false
    if item.selector == "#searchInput" {
      set isSearch = true
    }
    if item.selector == "#simpleSearch" {
      set isSearch = true
    }
    
    if isSearch == true {
      set searchInput = item
      print "Found search input: " + item.description
    }
  }
}

# Find search button
loop item in actionables {
  if item.type == "button" {
    if item.selector == "#searchButton" {
      set searchButton = item
      print "Found search button: " + item.description
    }
  }
}

# Interact with search if found
if searchInput != null {
  print "\nClicking on search input..."
  call browser_click {
    tabId: tab.id,
    selector: searchInput.selector
  }
  
  print "Typing search query..."
  call browser_type {
    tabId: tab.id,
    selector: searchInput.selector,
    text: "Artificial Intelligence",
    clearFirst: true
  }
  
  if searchButton != null {
    print "Clicking search button..."
    call browser_click {
      tabId: tab.id,
      selector: searchButton.selector
    }
    
    # Wait for navigation
    wait 3
    
    # Get actionables on the new page
    call browser_get_actionables {
      tabId: tab.id
    } -> newActionables
    
    print "\nAfter searching, found " + str(len(newActionables)) + " actionable elements on the AI article page"
  }
}

print "\n✓ Example completed!"
