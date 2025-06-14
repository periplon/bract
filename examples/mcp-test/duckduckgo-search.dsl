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

# Extract all links from the page (get href attribute to see actual URLs)
call browser_extract_content {
  tabId: tab.id,
  selector: "a[href]",
  contentType: "attribute",
  attribute: "href"
} -> allLinks1

# Extract all h2 elements
call browser_extract_content {
  tabId: tab.id,
  selector: "h2"
} -> allH2s1

# Extract all h3 elements (sometimes results use h3)
call browser_extract_content {
  tabId: tab.id,
  selector: "h3"
} -> allH3s1

# Print first search results
print "\nSearch Results for 'browser automation':"
print "Found " + str(len(allLinks1)) + " total links on page"
print "Found " + str(len(allH2s1)) + " h2 elements"
print "Found " + str(len(allH3s1)) + " h3 elements"

# Show all h2 elements
if len(allH2s1) > 0 {
  print "\nAll H2 elements:"
  set counter1 = 0
  loop h2 in allH2s1 {
    set counter1 = counter1 + 1
    print str(counter1) + ". " + h2
  }
}

# Show all h3 elements
if len(allH3s1) > 0 {
  print "\nAll H3 elements:"
  set counter1 = 0
  loop h3 in allH3s1 {
    set counter1 = counter1 + 1
    print str(counter1) + ". " + h3
  }
}

# Show sample of links that are likely search results
print "\nSample of links (showing first 20):"
set counter1 = 0
set maxLinks = 20
loop link in allLinks1 {
  if counter1 < maxLinks {
    # Filter for external links (actual search results)
    if link != "/" && link != "#" && link != "" {
      # Check if it's an external link (contains http and not duckduckgo)
      set isExternal = false
      if len(link) > 7 {
        if link != "https://duckduckgo.com" && link != "https://duckduckgo.com/" {
          set isExternal = true
        }
      }
      
      if isExternal {
        set counter1 = counter1 + 1
        print str(counter1) + ". " + link
      }
    }
  }
}

# Clear search box for second search
call browser_type {
  tabId: tab.id,
  selector: "input[name='q']",
  text: ""
}

# Second search: "web scraping tools"
print "\n\nSearching for 'web scraping tools'..."
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
wait 3

# Extract all content for second search
call browser_extract_content {
  tabId: tab.id,
  selector: "a[href]",
  contentType: "attribute",
  attribute: "href"
} -> allLinks2

call browser_extract_content {
  tabId: tab.id,
  selector: "h2"
} -> allH2s2

call browser_extract_content {
  tabId: tab.id,
  selector: "h3"
} -> allH3s2

print "\nSearch Results for 'web scraping tools':"
print "Found " + str(len(allLinks2)) + " total links on page"
print "Found " + str(len(allH2s2)) + " h2 elements"
print "Found " + str(len(allH3s2)) + " h3 elements"

# Show sample of links for second search
print "\nSample of links (showing first 10):"
set counter2 = 0
set maxLinks2 = 10
loop link in allLinks2 {
  if counter2 < maxLinks2 {
    # Filter for external links (actual search results)
    if link != "/" && link != "#" && link != "" {
      # Check if it's an external link (contains http and not duckduckgo)
      set isExternal = false
      if len(link) > 7 {
        if link != "https://duckduckgo.com" && link != "https://duckduckgo.com/" {
          set isExternal = true
        }
      }
      
      if isExternal {
        set counter2 = counter2 + 1
        print str(counter2) + ". " + link
      }
    }
  }
}

# Verify we got meaningful results
assert len(allLinks1) > 10, "Too few links found for browser automation"
assert len(allLinks2) > 10, "Too few links found for web scraping tools"

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ DuckDuckGo search test completed successfully"