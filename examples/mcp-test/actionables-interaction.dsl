# Actionables Interaction Example
# Shows how to discover and interact with page elements using browser_get_actionables

# Connect to the MCP browser server
connect "../../bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {
  timeout: 30
}
print "Browser connected"

# Navigate to example.com
call browser_navigate {
  url: "https://example.com",
  waitUntilLoad: true
}
print "Navigated to example.com"

# Get all actionable elements
print "\nAnalyzing page for actionable elements..."
call browser_get_actionables -> actionables

print "Found " + str(len(actionables)) + " actionable elements"

# Find all links and display them
print "\n=== LINKS ==="
set linkCount = 0
set firstLink = null
loop item in actionables {
  if item.type == "link" {
    print "[" + str(item.labelNumber) + "] " + item.description
    if linkCount == 0 {
      set firstLink = item
    }
    set linkCount = linkCount + 1
  }
}
print "Total links: " + str(linkCount)

# Find all buttons
print "\n=== BUTTONS ==="
set buttonCount = 0
loop item in actionables {
  if item.type == "button" {
    print "[" + str(item.labelNumber) + "] " + item.description
    set buttonCount = buttonCount + 1
  }
}
print "Total buttons: " + str(buttonCount)

# Find all input fields
print "\n=== INPUT FIELDS ==="
set inputCount = 0
loop item in actionables {
  if item.type == "input" {
    print "[" + str(item.labelNumber) + "] " + item.description + " (selector: " + item.selector + ")"
    set inputCount = inputCount + 1
  }
}
print "Total inputs: " + str(inputCount)

# Example: Click on the first link if available
if firstLink != null {
  print "\n=> Clicking on first link: " + firstLink.description
  call browser_click {
    selector: firstLink.selector
  }
  
  # Wait for navigation
  wait 2
  
  # Get actionables on the new page
  call browser_get_actionables -> newActionables
  print "New page has " + str(len(newActionables)) + " actionable elements"
  
  # Navigate back
  call browser_navigate {
    url: "https://example.com"
  }
  wait 1
}

# Example: Find and interact with a form if available
print "\n=> Looking for form elements..."
set formInput = null
set submitButton = null

# Find first input field
loop item in actionables {
  if item.type == "input" {
    set formInput = item
  }
}

# Find submit button
loop item in actionables {
  if item.type == "button" {
    set submitButton = item
  }
}

if formInput != null {
  print "Found input field: " + formInput.description
  
  # Click and type in the input
  call browser_click {
    selector: formInput.selector
  }
  
  call browser_type {
    selector: formInput.selector,
    text: "test input value",
    clearFirst: true
  }
  
  print "Typed test value in input field"
  
  if submitButton != null {
    print "Found button: " + submitButton.description
    # Note: Not clicking submit to avoid navigating away
  }
}

# Summary
print "\n=== Summary ==="
print "Page: example.com"
print "Total actionable elements: " + str(len(actionables))
print "- Links: " + str(linkCount)
print "- Buttons: " + str(buttonCount)
print "- Inputs: " + str(inputCount)

print "\n✓ Interaction example completed!"