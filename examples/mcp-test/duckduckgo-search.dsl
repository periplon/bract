# DuckDuckGo Search Test
# Demonstrates search functionality using DuckDuckGo

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension connection
call browser_wait_for_connection {timeout: 5}

# Create tab and navigate to DuckDuckGo
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://duckduckgo.com"
}

# Wait for search input to be available
call browser_wait_for_element {
  tabId: tab.id,
  selector: "input[name='q']"
}

# First search: "browser automation"
print "Searching for 'browser automation'..."
call browser_type {
  tabId: tab.id,
  selector: "input[name='q']",
  text: "browser automation"
}

# Submit search form
call browser_click {
  tabId: tab.id,
  selector: "button[type='submit']"
}

# Wait for results to load
wait 3

# Extract search results - using generic selectors
call browser_extract_content {
  tabId: tab.id,
  selector: "h2"
} -> titles1

# Get page content to verify results
call browser_extract_content {
  tabId: tab.id,
  selector: "body"
} -> pageContent1

# Check if we have content
if len(pageContent1) > 0 {
  print "Page content length: " + str(len(pageContent1[0]))
}

# Print first search results
print "\nSearch Results for 'browser automation':"
print "Found " + str(len(titles1)) + " h2 elements on results page"

if len(titles1) >= 1 {
  print "Sample result: " + titles1[0]
}

# Clear search box for second search
call browser_type {
  tabId: tab.id,
  selector: "input[name='q']",
  text: ""
}

# Second search: "web scraping tools"
print "\nSearching for 'web scraping tools'..."
call browser_type {
  tabId: tab.id,
  selector: "input[name='q']",
  text: "web scraping tools"
}

# Submit search form
call browser_click {
  tabId: tab.id,
  selector: "button[type='submit']"
}

# Wait for results to load
wait 2
call browser_wait_for_element {
  tabId: tab.id,
  selector: "[data-result]"
}

# Extract second search results - using h2 tags
call browser_extract_content {
  tabId: tab.id,
  selector: "h2"
} -> titles2

print "\nFound " + str(len(titles2)) + " results for 'web scraping tools'"

# Verify we got meaningful results
assert len(titles1) > 0, "No results found for browser automation"
assert len(titles2) > 0, "No results found for web scraping tools"

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ DuckDuckGo search test completed successfully"