# Form Interaction Test
# Tests browser form filling and interaction capabilities

# Connect to the MCP browser server
connect "./mcp-browser-server"

# Define reusable automation for form testing
define test_form(url, username, password) {
  # Create tab and navigate
  call create_tab -> tab
  call navigate {tabId: tab.id, url: url}
  
  # Wait for form to be present
  call wait_for_element {
    tabId: tab.id,
    selector: "form"
  }
  
  # Fill username field
  call type {
    tabId: tab.id,
    selector: "input[name='username']",
    text: username
  }
  
  # Fill password field
  call type {
    tabId: tab.id,
    selector: "input[name='password']",
    text: password
  }
  
  # Click submit button
  call click {
    tabId: tab.id,
    selector: "button[type='submit']"
  }
  
  # Wait for response
  wait 2
  
  # Check result
  call extract_content {
    tabId: tab.id,
    selector: ".message, .error, .success"
  } -> result
  
  print "Form submission result: " + result.text
  
  # Clean up
  call close_tab {tabId: tab.id}
}

# Run the form test automation
run test_form("https://example.com/login", "testuser", "testpass123")

print "✓ Form interaction test completed"