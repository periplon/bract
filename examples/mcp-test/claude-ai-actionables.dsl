# Claude AI Actionables Example
# Gets all actionable elements from claude.ai and lists them categorized by type

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
  timeout: 30
} -> connection_result
print "Browser connection established"

# Create a tab and navigate to Claude AI
call browser_create_tab {
  url: "https://claude.ai",
  active: true
} -> tab
print "Created tab with ID: " + str(tab.id)

# Wait for page to load completely
wait 3

# Get all actionable elements
print "\nAnalyzing claude.ai for actionable elements..."
call browser_get_actionables {
  tabId: tab.id
} -> actionables

print "\n============================================"
print "CLAUDE AI ACTIONABLES REPORT"
print "============================================"
print "Total actionable elements found: " + str(len(actionables))

# Initialize counters
set links = []
set buttons = []
set inputs = []
set textareas = []
set selects = []
set checkboxes = []
set radios = []
set others = []

# Categorize elements
loop item in actionables {
  if item.type == "a" {
    set links = links + [item]
  }
  if item.type == "button" {
    set buttons = buttons + [item]
  }
  if item.type == "input" {
    set inputs = inputs + [item]
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
  if item.type == "input[type=\"radio\"]" {
    set radios = radios + [item]
  }
  # Catch-all for other types
  if item.type != "a" && item.type != "button" && item.type != "input" && item.type != "textarea" && item.type != "select" && item.type != "input[type=\"checkbox\"]" && item.type != "input[type=\"radio\"]" {
    set others = others + [item]
  }
}

# Display summary by type
print "\n=== SUMMARY BY TYPE ==="
print "Links: " + str(len(links))
print "Buttons: " + str(len(buttons))
print "Input fields: " + str(len(inputs))
print "Textareas: " + str(len(textareas))
print "Selects: " + str(len(selects))
print "Checkboxes: " + str(len(checkboxes))
print "Radio buttons: " + str(len(radios))
print "Other elements: " + str(len(others))

# Display detailed list for each type
if len(links) > 0 {
  print "\n=== LINKS (" + str(len(links)) + ") ==="
  set count = 0
  loop link in links {
    if count < 10 {  # Show first 10
      print "[" + str(link.labelNumber) + "] " + link.description
      print "    Selector: " + link.selector
      set count = count + 1
    }
  }
  if len(links) > 10 {
    print "... and " + str(len(links) - 10) + " more links"
  }
}

if len(buttons) > 0 {
  print "\n=== BUTTONS (" + str(len(buttons)) + ") ==="
  set count = 0
  loop button in buttons {
    if count < 10 {  # Show first 10
      print "[" + str(button.labelNumber) + "] " + button.description
      print "    Selector: " + button.selector
      set count = count + 1
    }
  }
  if len(buttons) > 10 {
    print "... and " + str(len(buttons) - 10) + " more buttons"
  }
}

if len(inputs) > 0 {
  print "\n=== INPUT FIELDS (" + str(len(inputs)) + ") ==="
  loop input in inputs {
    print "[" + str(input.labelNumber) + "] " + input.description
    print "    Selector: " + input.selector
  }
}

if len(textareas) > 0 {
  print "\n=== TEXTAREAS (" + str(len(textareas)) + ") ==="
  loop textarea in textareas {
    print "[" + str(textarea.labelNumber) + "] " + textarea.description
    print "    Selector: " + textarea.selector
  }
}

if len(selects) > 0 {
  print "\n=== SELECT DROPDOWNS (" + str(len(selects)) + ") ==="
  loop select in selects {
    print "[" + str(select.labelNumber) + "] " + select.description
    print "    Selector: " + select.selector
  }
}

if len(checkboxes) > 0 {
  print "\n=== CHECKBOXES (" + str(len(checkboxes)) + ") ==="
  loop checkbox in checkboxes {
    print "[" + str(checkbox.labelNumber) + "] " + checkbox.description
    print "    Selector: " + checkbox.selector
  }
}

if len(radios) > 0 {
  print "\n=== RADIO BUTTONS (" + str(len(radios)) + ") ==="
  loop radio in radios {
    print "[" + str(radio.labelNumber) + "] " + radio.description
    print "    Selector: " + radio.selector
  }
}

if len(others) > 0 {
  print "\n=== OTHER ELEMENTS (" + str(len(others)) + ") ==="
  set count = 0
  loop other in others {
    if count < 5 {  # Show first 5
      print "[" + str(other.labelNumber) + "] Type: " + other.type + " - " + other.description
      print "    Selector: " + other.selector
      set count = count + 1
    }
  }
  if len(others) > 5 {
    print "... and " + str(len(others) - 5) + " more elements"
  }
}

# Look for specific elements that might be useful
print "\n=== NOTABLE ELEMENTS ==="

# Find the main chat input if available
set chatInput = null
loop item in actionables {
  if item.type == "textarea" || item.type == "input" {
    # Check if description mentions chat, message, or prompt
    if item.description == "Enter your message" || item.description == "Type a message" || item.description == "Chat input" || item.description == "Message input" {
      set chatInput = item
    }
  }
}

if chatInput != null {
  print "✓ Found chat input: " + chatInput.description
  print "  Selector: " + chatInput.selector
} else {
  print "✗ No obvious chat input field found"
}

# Find the send/submit button
set sendButton = null
loop item in actionables {
  if item.type == "button" {
    if item.description == "Send" || item.description == "Submit" || item.description == "Send message" {
      set sendButton = item
    }
  }
}

if sendButton != null {
  print "✓ Found send button: " + sendButton.description
  print "  Selector: " + sendButton.selector
} else {
  print "✗ No obvious send button found"
}

# Find new chat button
set newChatButton = null
loop item in actionables {
  if item.type == "button" || item.type == "a" {
    if item.description == "New chat" || item.description == "Start new chat" || item.description == "New conversation" {
      set newChatButton = item
    }
  }
}

if newChatButton != null {
  print "✓ Found new chat button: " + newChatButton.description
  print "  Selector: " + newChatButton.selector
} else {
  print "✗ No obvious new chat button found"
}

print "\n============================================"
print "✓ Claude AI actionables analysis complete!"
print "============================================"