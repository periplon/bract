# Get Accessibility Snapshot of Current Tab
# Simple example showing how to analyze the accessibility tree of the active tab

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

print "=== Accessibility Snapshot of Current Tab ==="

# Get accessibility snapshot of the current active tab (tabId defaults to 0 = active tab)
call browser_get_accessibility_snapshot {} -> snapshot

# Display basic information
print "\nPage Accessibility Information:"
print "Role: " + snapshot.role
if snapshot.name {
  print "Name: " + snapshot.name
}

# Count different types of elements
elementCounts = {}

function countElements(node) {
  role = node.role || "unknown"
  if elementCounts[role] {
    elementCounts[role] = elementCounts[role] + 1
  } else {
    elementCounts[role] = 1
  }
  
  if node.children {
    for child in node.children {
      countElements(child)
    }
  }
}

countElements(snapshot)

print "\nElement Types Found:"
for role in Object.keys(elementCounts).sort() {
  print "  " + role + ": " + elementCounts[role]
}

# Show page structure
print "\nPage Structure (first 2 levels):"
function showStructure(node, indent, maxDepth) {
  if maxDepth <= 0 {
    return
  }
  
  prefix = "  ".repeat(indent)
  nodeInfo = node.role
  if node.name {
    nodeInfo = nodeInfo + ' "' + node.name + '"'
  }
  
  print prefix + "- " + nodeInfo
  
  if node.children && node.children.length > 0 {
    # Show up to 5 children per level
    childrenToShow = Math.min(node.children.length, 5)
    for i in range(0, childrenToShow) {
      showStructure(node.children[i], indent + 1, maxDepth - 1)
    }
    
    if node.children.length > 5 {
      print prefix + "  ... and " + (node.children.length - 5) + " more"
    }
  }
}

showStructure(snapshot, 0, 2)

print "\n✓ Accessibility analysis complete!"