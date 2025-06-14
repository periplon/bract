# Form Interaction Test
# Tests browser form filling and interaction capabilities

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension connection
call browser_wait_for_connection {timeout: 5}

# Define reusable automation for form testing
define test_form(url, username, password) {
  # Create tab and navigate
  call browser_create_tab -> tab
  call browser_navigate {tabId: tab.id, url: url}
  
  # Wait for form to be present
  call browser_wait_for_element {
    tabId: tab.id,
    selector: "form"
  }
  
  # Fill username field
  call browser_type {
    tabId: tab.id,
    selector: "input[name='username']",
    text: username
  }
  
  # Fill password field
  call browser_type {
    tabId: tab.id,
    selector: "input[name='password']",
    text: password
  }
  
  # Click submit button
  call browser_click {
    tabId: tab.id,
    selector: "button[type='submit']"
  }
  
  # Wait for response
  wait 2
  
  # Check result
  call browser_extract_content {
    tabId: tab.id,
    selector: ".message, .error, .success, body"
  } -> result
  
  if len(result) > 0 {
    print "Form submission result: " + result[0]
  } else {
    print "Form submitted (no message found)"
  }
  
  # Clean up
  call browser_close_tab {tabId: tab.id}
}

# Run the form test automation
run test_form("https://example.com/login", "testuser", "testpass123")

print "✓ Form interaction test completed"