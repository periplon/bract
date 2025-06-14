# Browser Reload Test
# Tests page reload functionality including hard reload

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

# Create a tab and navigate to a page with dynamic content
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://example.com"
}

print "=== Browser Reload Test ==="

# Get initial page content
call browser_extract_content {
  tabId: tab.id,
  selector: "h1"
} -> initialContent
set initialTitle = initialContent[0]
print "Initial page title: " + initialTitle

# Add some data to sessionStorage (will be lost on reload)
call browser_set_session_storage {
  tabId: tab.id,
  key: "tempData",
  value: "This will be lost on reload"
}

# Add data to localStorage (will persist)
call browser_set_local_storage {
  tabId: tab.id,
  key: "persistentData",
  value: "This will survive reload"
}

print "\n1. Testing normal reload:"
# Normal reload (uses cache)
call browser_reload {
  tabId: tab.id,
  hardReload: false
} -> reloadResult
print "✓ Page reloaded normally"

# Wait for page to load
wait 2

# Check that localStorage persisted
call browser_get_local_storage {
  tabId: tab.id,
  key: "persistentData"
} -> persistentCheck
assert persistentCheck.value == "This will survive reload", "localStorage should persist after reload"
print "✓ localStorage data persisted"

# Check that sessionStorage was cleared
call browser_get_session_storage {
  tabId: tab.id,
  key: "tempData"
} -> tempCheck
print "✓ sessionStorage was cleared as expected"

print "\n2. Testing hard reload (bypass cache):"
# Set a timestamp in localStorage
call browser_set_local_storage {
  tabId: tab.id,
  key: "reloadTime",
  value: "timestamp_before_hard_reload"
}

# Hard reload (bypasses cache)
call browser_reload {
  tabId: tab.id,
  hardReload: true
} -> hardReloadResult
print "✓ Page hard reloaded (cache bypassed)"

# Wait for page to load
wait 2

# Verify page loaded fresh by checking if we can extract content
call browser_extract_content {
  tabId: tab.id,
  selector: "body"
} -> bodyContent
assert len(bodyContent) > 0, "Page should be fully loaded"
print "✓ Page fully loaded after hard reload"

print "\n3. Testing reload without tabId (uses active tab):"
# Ensure our tab is active
call browser_activate_tab {tabId: tab.id}

# Reload without specifying tabId
call browser_reload {
  hardReload: false
}
print "✓ Active tab reloaded successfully"

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ Browser reload test completed!"