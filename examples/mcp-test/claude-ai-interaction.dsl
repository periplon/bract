# DSL script to navigate to Claude AI and type in the contenteditable div
# This script opens a new tab, navigates to the Claude AI new chat page,
# and types "who am I" in the main input area

# Connect to the browser MCP server
connect "./bin/mcp-browser-server"

# Wait for connection to be established
call browser_wait_for_connection {
  timeout: 5
}

# Create a new browser tab
call browser_create_tab -> tab

# Navigate to Claude AI new chat page
call browser_navigate {
  tabId: tab.id,
  url: "https://claude.ai/new"
}

# Wait for the page to load completely
wait 3

# Wait for the contenteditable div to be available
call browser_wait_for_element {
  tabId: tab.id,
  selector: "div[contenteditable='true']",
  timeout: 10000
}

# Additional wait to ensure the element is fully interactive
wait 1

# Type "who am I" in the contenteditable div
call browser_type {
  tabId: tab.id,
  selector: "div[contenteditable='true']",
  text: "who am I"
}

# Optional: Take a screenshot to verify the input
call browser_screenshot {
  tabId: tab.id,
  fullPage: false
} -> screenshot

# Print confirmation
print "Successfully typed 'who am I' in the Claude AI input field"

# Keep the tab open for manual verification
# Uncomment the line below to close the tab after completion
# call browser_close_tab {tabId: tab.id}