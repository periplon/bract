# ChatGPT Actionables Example
# Gets all actionable elements from chatgpt.com and lists them categorized by type

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
  timeout: 30
} -> connection_result
print "Browser connection established"

# Create a tab and navigate to ChatGPT
call browser_create_tab {
  url: "https://chatgpt.com",
  active: true
} -> tab
print "Created tab with ID: " + str(tab.id)

# Wait for page to load completely
wait 3

# Get all actionable elements
print "\nAnalyzing chatgpt.com for actionable elements..."
call browser_get_actionables {
  tabId: tab.id
} -> actionables

print "\n============================================"
print "CHATGPT ACTIONABLES REPORT"
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
    if count < 15 {  # Show first 15 for ChatGPT as it has many buttons
      print "[" + str(button.labelNumber) + "] " + button.description
      print "    Selector: " + button.selector
      set count = count + 1
    }
  }
  if len(buttons) > 15 {
    print "... and " + str(len(buttons) - 15) + " more buttons"
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

# Look for specific ChatGPT UI elements
print "\n=== NOTABLE CHATGPT ELEMENTS ==="

# Find the main chat input/prompt field
set chatInput = null
loop item in actionables {
  if item.type == "textarea" || item.type == "input" {
    # Check for common ChatGPT input field descriptions
    if item.description == "Message ChatGPT" || item.description == "Send a message" || item.description == "Type your message" || item.description == "Message" || item.description == "Send a message..." {
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
    # Look for send button (might be an icon button)
    if item.description == "Send" || item.description == "Send message" || item.description == "Submit" || item.description == "Send prompt" {
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
    if item.description == "New chat" || item.description == "New Chat" || item.description == "Start new chat" || item.description == "Create new chat" {
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

# Find model selector if available
set modelSelector = null
loop item in actionables {
  if item.type == "button" || item.type == "select" {
    # Look for GPT model selector
    if item.description == "GPT-4" || item.description == "GPT-3.5" || item.description == "Model" || item.description == "Select model" {
      set modelSelector = item
    }
  }
}

if modelSelector != null {
  print "✓ Found model selector: " + modelSelector.description
  print "  Selector: " + modelSelector.selector
} else {
  print "✗ No model selector found"
}

# Find settings/menu buttons
set settingsButton = null
loop item in actionables {
  if item.type == "button" || item.type == "a" {
    if item.description == "Settings" || item.description == "Menu" || item.description == "User menu" || item.description == "Options" {
      set settingsButton = item
    }
  }
}

if settingsButton != null {
  print "✓ Found settings/menu button: " + settingsButton.description
  print "  Selector: " + settingsButton.selector
} else {
  print "✗ No settings/menu button found"
}

# Look for conversation history items
print "\n=== CONVERSATION ELEMENTS ==="
set conversationCount = 0
loop item in actionables {
  if item.type == "a" || item.type == "button" {
    # Check if it looks like a conversation item
    if item.description == "Previous conversation" || item.description == "Chat history" {
      set conversationCount = conversationCount + 1
    }
  }
}

if conversationCount > 0 {
  print "✓ Found " + str(conversationCount) + " conversation history items"
} else {
  print "✗ No conversation history items found"
}

print "\n============================================"
print "✓ ChatGPT actionables analysis complete!"
print "============================================"