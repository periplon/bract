# Advanced Browser Automation
# Demonstrates complex automation patterns and reusable components

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

# Define helper functions
define wait_for_element_helper(tabId, selector) {
  call browser_wait_for_element {
    tabId: tabId,
    selector: selector,
    timeout: 5
  }
}

define fill_input_helper(tabId, selector, text) {
  call browser_type {
    tabId: tabId,
    selector: selector,
    text: text
  }
}

define click_element_helper(tabId, selector) {
  call browser_click {
    tabId: tabId,
    selector: selector
  }
}

# Define a test scenario using example.com
define test_navigation_functionality(testUrl) {
  print "Testing navigation to: " + testUrl
  
  # Create tab
  call browser_create_tab -> tab
  
  # Navigate to URL
  call browser_navigate {
    tabId: tab.id,
    url: testUrl
  }
  
  # Wait for page to load
  wait 2
  
  # Extract page content to verify navigation worked
  call browser_extract_content {
    tabId: tab.id,
    selector: "h1",
    contentType: "text"
  } -> headings
  
  set headingCount = len(headings)
  print "Found " + str(headingCount) + " h1 elements"
  assert headingCount > 0, "No h1 elements found on page"
  
  if headingCount > 0 {
    print "First heading: " + headings[0]
  }
  
  # Take screenshot as additional verification
  call browser_screenshot {
    tabId: tab.id
  } -> screenshot
  assert screenshot.dataUrl != null, "Screenshot failed"
  print "✓ Screenshot captured"
  
  # Clean up
  call browser_close_tab {tabId: tab.id}
  
  print "✓ Navigation test passed for: " + testUrl
}

# Define batch testing
define run_test_suite(testCases) {
  set passed = 0
  set failed = 0
  
  print "Running test suite with " + str(len(testCases)) + " test cases\n"
  
  loop testCase in testCases {
    print "─────────────────────────"
    # Run test in try-catch pattern (using if with error checking)
    set error = false
    
    if testCase.type == "navigation" {
      run test_navigation_functionality(testCase.data)
    }
    
    if !error {
      set passed = passed + 1
    } else {
      set failed = failed + 1
    }
  }
  
  print "\n===== Test Suite Results ====="
  print "Passed: " + str(passed)
  print "Failed: " + str(failed)
  print "Total:  " + str(len(testCases))
  print "Success Rate: " + str((passed * 100) / len(testCases)) + "%"
}

# Run the test suite
set testCases = [{type: "navigation", data: "https://example.com"}, {type: "navigation", data: "https://example.org"}, {type: "navigation", data: "https://example.net"}]

run run_test_suite(testCases)

print "\n✓ Advanced automation demo completed"