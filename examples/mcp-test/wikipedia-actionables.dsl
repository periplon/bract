# Wikipedia Actionables Example
# This example demonstrates how to use the browser_get_actionables tool
# to discover all interactive elements on a Wikipedia page

# First, wait for the browser extension to connect
tool browser_wait_for_connection

# Create a new tab and navigate to Wikipedia
result = tool browser_create_tab {
    "url": "https://en.wikipedia.org/wiki/Main_Page",
    "active": true
}

# Extract the tab ID from the result
tabId = result.id

# Wait for the page to fully load
tool browser_wait_for_element {
    "tabId": tabId,
    "selector": "#mp-welcome",
    "timeout": 10000
}

# Get all actionable elements on the Wikipedia main page
actionables = tool browser_get_actionables {
    "tabId": tabId
}

# Print the total number of actionable elements found
print "Found " + len(actionables) + " actionable elements on Wikipedia's main page"
print "---"

# Display the first 20 actionable elements with their details
counter = 0
for actionable in actionables {
    if counter < 20 {
        print "Element #" + actionable.labelNumber + ":"
        print "  Description: " + actionable.description
        print "  Type: " + actionable.type
        print "  Selector: " + actionable.selector
        print ""
        counter = counter + 1
    }
}

# Find and display specific types of elements
print "=== Links ==="
linkCount = 0
for actionable in actionables {
    if actionable.type == "link" && linkCount < 10 {
        print "- " + actionable.description + " (" + actionable.selector + ")"
        linkCount = linkCount + 1
    }
}

print ""
print "=== Buttons ==="
buttonCount = 0
for actionable in actionables {
    if actionable.type == "button" && buttonCount < 10 {
        print "- " + actionable.description + " (" + actionable.selector + ")"
        buttonCount = buttonCount + 1
    }
}

print ""
print "=== Input Fields ==="
inputCount = 0
for actionable in actionables {
    if actionable.type == "input" && inputCount < 10 {
        print "- " + actionable.description + " (" + actionable.selector + ")"
        inputCount = inputCount + 1
    }
}

# Example: Click on the search input field if found
searchFound = false
for actionable in actionables {
    if actionable.type == "input" && (actionable.description contains "Search" || actionable.selector contains "search") {
        print ""
        print "Found search input! Clicking on it..."
        tool browser_click {
            "tabId": tabId,
            "selector": actionable.selector
        }
        
        # Type a search query
        tool browser_type {
            "tabId": tabId,
            "selector": actionable.selector,
            "text": "Artificial Intelligence",
            "clearFirst": true
        }
        
        searchFound = true
        break
    }
}

if searchFound {
    print "Typed 'Artificial Intelligence' in the search box"
    
    # Find and click the search button
    for actionable in actionables {
        if actionable.type == "button" && (actionable.description contains "Search" || actionable.description contains "Go") {
            print "Clicking search button..."
            tool browser_click {
                "tabId": tabId,
                "selector": actionable.selector
            }
            
            # Wait for navigation
            wait 3
            
            # Get actionables on the new page
            newActionables = tool browser_get_actionables {
                "tabId": tabId
            }
            
            print ""
            print "After searching, found " + len(newActionables) + " actionable elements on the AI article page"
            break
        }
    }
}

print ""
print "Example completed!"