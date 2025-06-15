# Accessibility Snapshot Example
# Demonstrates how to get the accessibility tree of a web page

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

# List existing tabs
call browser_list_tabs -> tabs

if tabs.length == 0 {
  # Create a new tab if none exist
  call browser_create_tab {
    url: "https://www.w3.org/WAI/ARIA/apg/patterns/",
    active: true
  } -> tab
} else {
  # Use the first tab
  tab = tabs[0]
  
  # Navigate to the ARIA patterns page
  call browser_navigate {
    tabId: tab.id,
    url: "https://www.w3.org/WAI/ARIA/apg/patterns/"
  }
}

print "=== Accessibility Snapshot Example ==="
print "Current tab: " + tab.title + " (" + tab.url + ")"

# Wait for page to load
call browser_wait_for_element {
  tabId: tab.id,
  selector: "main",
  timeout: 5000
}

# Get the full accessibility snapshot
print "\n1. Full Accessibility Snapshot (interesting nodes only):"
call browser_get_accessibility_snapshot {
  tabId: tab.id,
  interestingOnly: true
} -> snapshot

print "Root role: " + snapshot.role
print "Page name: " + snapshot.name
print "Number of direct children: " + snapshot.children.length

# Get a more detailed snapshot (all nodes)
print "\n2. Detailed Accessibility Snapshot (all nodes):"
call browser_get_accessibility_snapshot {
  tabId: tab.id,
  interestingOnly: false
} -> detailedSnapshot

print "Total accessible elements found (including all nodes)"

# Get accessibility snapshot of a specific region
print "\n3. Accessibility Snapshot of Main Content:"
call browser_get_accessibility_snapshot {
  tabId: tab.id,
  interestingOnly: true,
  root: "main"
} -> mainSnapshot

print "Main content role: " + mainSnapshot.role
if mainSnapshot.name {
  print "Main content name: " + mainSnapshot.name
}

# Find all headings in the accessibility tree
print "\n4. Finding Headings in Accessibility Tree:"
headings = []

function findHeadings(node) {
  if node.role == "heading" {
    headings.push({
      name: node.name,
      level: node.level
    })
  }
  
  if node.children {
    for child in node.children {
      findHeadings(child)
    }
  }
}

findHeadings(snapshot)

print "Found " + headings.length + " headings:"
for heading in headings {
  print "  Level " + heading.level + ": " + heading.name
}

# Find all interactive elements
print "\n5. Finding Interactive Elements:"
interactiveElements = []

function findInteractive(node) {
  interactiveRoles = ["button", "link", "textbox", "checkbox", "radio", "combobox", "menuitem"]
  
  if interactiveRoles.includes(node.role) {
    interactiveElements.push({
      role: node.role,
      name: node.name || "(unnamed)",
      disabled: node.disabled || false
    })
  }
  
  if node.children {
    for child in node.children {
      findInteractive(child)
    }
  }
}

findInteractive(snapshot)

print "Found " + interactiveElements.length + " interactive elements:"
for elem in interactiveElements.slice(0, 10) {  # Show first 10
  status = elem.disabled ? " (disabled)" : ""
  print "  " + elem.role + ": " + elem.name + status
}

if interactiveElements.length > 10 {
  print "  ... and " + (interactiveElements.length - 10) + " more"
}

print "\n✓ Accessibility snapshot example completed!"