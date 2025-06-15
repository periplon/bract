# Actionables Interaction Example
# Shows how to discover and interact with page elements using browser_get_actionables

# Connect and navigate
tool browser_wait_for_connection
tool browser_navigate {
    "url": "https://example.com",
    "waitUntilLoad": true
}

# Get all actionable elements
print "Analyzing page for actionable elements..."
actionables = tool browser_get_actionables

print "Found " + len(actionables) + " actionable elements"

# Find all links and display them
print "\n=== LINKS ==="
links = []
for item in actionables {
    if item.type == "link" {
        links.append(item)
        print "[" + str(item.labelNumber) + "] " + item.description
    }
}

# Find all buttons
print "\n=== BUTTONS ==="
buttons = []
for item in actionables {
    if item.type == "button" {
        buttons.append(item)
        print "[" + str(item.labelNumber) + "] " + item.description
    }
}

# Find all input fields
print "\n=== INPUT FIELDS ==="
inputs = []
for item in actionables {
    if item.type == "input" {
        inputs.append(item)
        print "[" + str(item.labelNumber) + "] " + item.description + " (selector: " + item.selector + ")"
    }
}

# Example: Click on the first link if available
if len(links) > 0 {
    firstLink = links[0]
    print "\n=> Clicking on first link: " + firstLink.description
    tool browser_click {
        "selector": firstLink.selector
    }
    
    # Wait for navigation
    wait 2
    
    # Get actionables on the new page
    newActionables = tool browser_get_actionables
    print "New page has " + len(newActionables) + " actionable elements"
}

# Example: Find and interact with a search box
print "\n=> Looking for search functionality..."
searchInput = null
for item in actionables {
    desc = item.description.lower()
    selector = item.selector.lower()
    if item.type == "input" && (desc contains "search" || selector contains "search") {
        searchInput = item
        break
    }
}

if searchInput != null {
    print "Found search input: " + searchInput.description
    
    # Click and type in the search box
    tool browser_click {
        "selector": searchInput.selector
    }
    
    tool browser_type {
        "selector": searchInput.selector,
        "text": "test search query",
        "clearFirst": true
    }
    
    print "Typed search query"
    
    # Look for a search button near the input
    searchButton = null
    for item in actionables {
        if item.type == "button" && (item.description.lower() contains "search" || item.description.lower() contains "submit") {
            searchButton = item
            break
        }
    }
    
    if searchButton != null {
        print "Clicking search button: " + searchButton.description
        tool browser_click {
            "selector": searchButton.selector
        }
    }
} else {
    print "No search input found on this page"
}

print "\nInteraction example completed!"