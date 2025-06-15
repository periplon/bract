# VWO Form Testing Page Example
# Navigates to VWO glossary page and clicks on a specific link

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension connection
call browser_wait_for_connection {timeout: 5}

# Create tab and navigate to VWO form testing glossary page
call browser_create_tab -> tab
print "Navigating to VWO form testing glossary page..."
call browser_navigate {
  tabId: tab.id,
  url: "https://vwo.com/glossary/form-testing/"
}

# Wait for page to load
wait 3

# Extract page title to verify we're on the right page
call browser_list_tabs -> tabs
set pageTitle = ""
loop t in tabs {
  if t.id == tab.id {
    set pageTitle = t.title
  }
}
print "Page title: " + pageTitle

# Find all links with the specified class
print "\nLooking for links with class 'pink-common-link D(ib)'..."
call browser_extract_content {
  tabId: tab.id,
  selector: "a.pink-common-link.D\\(ib\\)"
} -> targetLinks

# Get all anchor elements to see their structure
call browser_extract_content {
  tabId: tab.id,
  selector: "a",
  contentType: "attribute",
  attribute: "href"
} -> allHrefs

print "Total anchors on page: " + str(len(allHrefs))

print "Found " + str(len(targetLinks)) + " links with the specified class"

# Display the links found
if len(targetLinks) > 0 {
  print "\nLinks found:"
  set counter = 0
  loop link in targetLinks {
    set counter = counter + 1
    print str(counter) + ". Text: " + link
  }
  
  # Click on the first link with the specified class
  print "\nClicking on the first link..."
  # The browser_click function doesn't return a value, so we shouldn't assign it
  call browser_click {
    tabId: tab.id,
    selector: "a.pink-common-link.D\\(ib\\)"
  }
  
  print "✓ Click command executed"
  
  # Wait for navigation
  wait 3
  
  # Get the new page title after clicking
  call browser_list_tabs -> newTabs
  set newPageTitle = ""
  loop t in newTabs {
    if t.id == tab.id {
      set newPageTitle = t.title
    }
  }
  print "New page title: " + newPageTitle
  
  # Get the current URL
  call browser_extract_content {
    tabId: tab.id,
    selector: "body"
  } -> bodyContent
  
  if len(bodyContent) > 0 {
    print "✓ Page loaded successfully"
  }
} else {
  print "✗ No links found with class 'pink-common-link D(ib)'"
  
  # Try alternative selectors
  print "\nTrying alternative selectors..."
  
  # Try without escaping parentheses
  call browser_extract_content {
    tabId: tab.id,
    selector: "a.pink-common-link"
  } -> alternativeLinks
  
  print "Found " + str(len(alternativeLinks)) + " links with class 'pink-common-link'"
  
  if len(alternativeLinks) > 0 {
    print "\nAlternative links found:"
    set counter = 0
    loop link in alternativeLinks {
      set counter = counter + 1
      print str(counter) + ". " + link
      if counter >= 5 {
        print "... (showing first 5)"
        set counter = -1
      }
    }
  }
}

# Take a screenshot for verification
print "\nTaking screenshot..."
call browser_screenshot {
  tabId: tab.id
} -> screenshot

if screenshot.dataUrl != null {
  print "✓ Screenshot captured"
} else {
  print "✗ Screenshot failed"
}

# Clean up
call browser_close_tab {
  tabId: tab.id
}

print "\n✓ VWO form testing example completed"