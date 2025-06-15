# Simple Get Actionables Example
# Demonstrates the basic usage of browser_get_actionables tool

# Connect to browser
tool browser_wait_for_connection

# Navigate to Wikipedia
tool browser_navigate {
    "url": "https://en.wikipedia.org",
    "waitUntilLoad": true
}

# Get all actionable elements
actionables = tool browser_get_actionables

# Display summary
print "Total actionable elements: " + len(actionables)

# Show first 10 elements
print "\nFirst 10 actionable elements:"
for i in range(0, min(10, len(actionables))) {
    item = actionables[i]
    print str(i+1) + ". [" + item.type + "] " + item.description
    print "   Selector: " + item.selector
}

# Count by type
types = {}
for item in actionables {
    if item.type in types {
        types[item.type] = types[item.type] + 1
    } else {
        types[item.type] = 1
    }
}

print "\nElements by type:"
for type, count in types {
    print "- " + type + ": " + str(count)
}