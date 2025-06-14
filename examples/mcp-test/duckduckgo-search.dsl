# DuckDuckGo Search Test
# Demonstrates search functionality using DuckDuckGo

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension connection
call browser_wait_for_connection {timeout: 5}

# Define reusable search automation
define search_duckduckgo(query, max_results = 5) {
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
  
  # Enter search query
  call browser_type {
    tabId: tab.id,
    selector: "input[name='q']",
    text: query
  }
  
  # Submit search form
  call browser_click {
    tabId: tab.id,
    selector: "button[type='submit']"
  }
  
  # Wait for results to load
  call browser_wait_for_element {
    tabId: tab.id,
    selector: "[data-result]"
  }
  
  # Extract search results
  call browser_extract_content {
    tabId: tab.id,
    selector: "[data-result] h2 a"
  } -> titles
  
  call browser_extract_content {
    tabId: tab.id,
    selector: "[data-result] .result__snippet"
  } -> snippets
  
  call browser_extract_content {
    tabId: tab.id,
    selector: "[data-result] .result__url"
  } -> urls
  
  # Format results
  set results = []
  loop i in range(min(len(titles), max_results)) {
    set result = {
      title: titles[i],
      snippet: snippets[i] if i < len(snippets) else "",
      url: urls[i] if i < len(urls) else ""
    }
    push results, result
  }
  
  # Clean up
  call browser_close_tab {tabId: tab.id}
  
  return results
}

# Perform searches
print "Searching for 'browser automation'..."
run search_duckduckgo("browser automation", 3) -> automation_results

print "Search Results:"
loop result in automation_results {
  print "Title: " + result.title
  print "URL: " + result.url
  print "Snippet: " + result.snippet
  print "---"
}

print "Searching for 'web scraping tools'..."
run search_duckduckgo("web scraping tools", 5) -> scraping_results

print "Found " + str(len(scraping_results)) + " results for web scraping tools"

# Verify we got meaningful results
assert len(automation_results) > 0, "No results found for browser automation"
assert len(scraping_results) > 0, "No results found for web scraping tools"

print "✓ DuckDuckGo search test completed successfully"