# Advanced Browser Automation
# Demonstrates complex automation patterns and reusable components

# Connect to the MCP browser server
connect "./mcp-browser-server"

# Define a page object pattern
define PageObject(tabId) {
  # Store tab ID for all operations
  set this_tab = tabId
  
  # Helper to wait for element
  define wait_for(selector) {
    call wait_for_element {
      tabId: this_tab,
      selector: selector,
      timeout: 5000
    }
  }
  
  # Helper to get element text
  define get_text(selector) {
    call extract_content {
      tabId: this_tab,
      selector: selector
    } -> content
    set result = content.text
  }
  
  # Helper to fill input
  define fill_input(selector, text) {
    call type {
      tabId: this_tab,
      selector: selector,
      text: text
    }
  }
  
  # Helper to click element
  define click_element(selector) {
    call click {
      tabId: this_tab,
      selector: selector
    }
  }
}

# Define a test scenario
define test_search_functionality(searchTerm) {
  print "Testing search with term: " + searchTerm
  
  # Create tab
  call create_tab -> tab
  
  # Initialize page object
  run PageObject(tab.id)
  
  # Navigate to search engine
  call navigate {
    tabId: tab.id,
    url: "https://duckduckgo.com"
  }
  
  # Wait for search box
  run wait_for("input[name='q']")
  
  # Perform search
  run fill_input("input[name='q']", searchTerm)
  run click_element("button[type='submit']")
  
  # Wait for results
  wait 2
  run wait_for(".results")
  
  # Count results
  call execute_script {
    tabId: tab.id,
    script: "document.querySelectorAll('.result').length"
  } -> resultCount
  
  print "Found " + str(resultCount.result) + " results"
  assert resultCount.result > 0, "No search results found"
  
  # Clean up
  call close_tab {tabId: tab.id}
  
  print "✓ Search test passed for: " + searchTerm
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
    
    if testCase.type == "search" {
      run test_search_functionality(testCase.data)
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
set testCases = [
  {type: "search", data: "MCP protocol"},
  {type: "search", data: "browser automation"},
  {type: "search", data: "test automation DSL"}
]

run run_test_suite(testCases)

print "\n✓ Advanced automation demo completed"