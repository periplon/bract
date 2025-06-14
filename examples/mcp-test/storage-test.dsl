# Storage Test
# Tests browser storage capabilities (cookies, localStorage, sessionStorage)

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
  timeout: 5
} -> connection_result
print "Browser connection status:"
print connection_result

# Create a test tab
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://example.com"
}

# Test Cookies
print "\n=== Testing Cookies ==="

# Set a cookie
call browser_set_cookie {
  name: "test_cookie",
  value: "test_value",
  domain: ".example.com",
  path: "/",
  secure: true,
  httpOnly: false
} -> cookie_result
print "Set cookie result:"
print cookie_result

# Get cookies
call browser_get_cookies {
  url: "https://example.com"
} -> cookies
print "All cookies for example.com:"
print cookies

# Get specific cookie
call browser_get_cookies {
  url: "https://example.com",
  name: "test_cookie"
} -> specific_cookie
print "Specific cookie 'test_cookie':"
print specific_cookie

# Delete specific cookie
call browser_delete_cookies {
  url: "https://example.com",
  name: "test_cookie"
} -> delete_result
print "Deleted cookie 'test_cookie'"

# Verify deletion
call browser_get_cookies {
  url: "https://example.com",
  name: "test_cookie"
} -> deleted_check
print "Cookies after deletion:"
print deleted_check

# Test LocalStorage
print "\n=== Testing LocalStorage ==="

# Set localStorage items
call browser_set_local_storage {
  tabId: tab.id,
  key: "test_key",
  value: "test_value"
} -> ls_set_result
print "Set localStorage['test_key'] = 'test_value'"

call browser_set_local_storage {
  tabId: tab.id,
  key: "another_key",
  value: "another_value"
} -> ls_set_result2
print "Set localStorage['another_key'] = 'another_value'"

# Get localStorage items
call browser_get_local_storage {
  tabId: tab.id,
  key: "test_key"
} -> ls_value
print "Retrieved localStorage['test_key']:"
print ls_value

call browser_get_local_storage {
  tabId: tab.id,
  key: "another_key"
} -> ls_value2
print "Retrieved localStorage['another_key']:"
print ls_value2

# Clear localStorage
call browser_clear_local_storage {
  tabId: tab.id
} -> ls_clear_result
print "Cleared all localStorage"

# Verify localStorage is cleared - getting a non-existent key may return empty string or error
call browser_get_local_storage {
  tabId: tab.id,
  key: "test_key"
} -> ls_cleared_check
print "localStorage['test_key'] after clear (should be empty):"
print ls_cleared_check

# Test SessionStorage
print "\n=== Testing SessionStorage ==="

# Set sessionStorage items
call browser_set_session_storage {
  tabId: tab.id,
  key: "session_key",
  value: "session_value"
} -> ss_set_result
print "Set sessionStorage['session_key'] = 'session_value'"

call browser_set_session_storage {
  tabId: tab.id,
  key: "temp_key",
  value: "temp_value"
} -> ss_set_result2
print "Set sessionStorage['temp_key'] = 'temp_value'"

# Get sessionStorage items
call browser_get_session_storage {
  tabId: tab.id,
  key: "session_key"
} -> ss_value
print "Retrieved sessionStorage['session_key']:"
print ss_value

call browser_get_session_storage {
  tabId: tab.id,
  key: "temp_key"
} -> ss_value2
print "Retrieved sessionStorage['temp_key']:"
print ss_value2

# Clear sessionStorage
call browser_clear_session_storage {
  tabId: tab.id
} -> ss_clear_result
print "Cleared all sessionStorage"

# Verify sessionStorage is cleared - getting a non-existent key may return empty string or error
call browser_get_session_storage {
  tabId: tab.id,
  key: "session_key"
} -> ss_cleared_check
print "sessionStorage['session_key'] after clear (should be empty):"
print ss_cleared_check

# Test storage persistence across navigation
print "\n=== Testing Storage Persistence ==="

# Set storage before navigation
call browser_set_local_storage {
  tabId: tab.id,
  key: "persist_key",
  value: "persist_value"
} -> persist_set
print "Set localStorage['persist_key'] before navigation"

call browser_set_session_storage {
  tabId: tab.id,
  key: "session_persist_key",
  value: "session_persist_value"
} -> session_persist_set
print "Set sessionStorage['session_persist_key'] before navigation"

# Navigate to another page
call browser_navigate {
  tabId: tab.id,
  url: "https://example.com/test"
}
print "Navigated to /test"

# Check localStorage (should persist)
call browser_get_local_storage {
  tabId: tab.id,
  key: "persist_key"
} -> persist_check
print "localStorage['persist_key'] after navigation:"
print persist_check

# Check sessionStorage (should persist for same origin)
call browser_get_session_storage {
  tabId: tab.id,
  key: "session_persist_key"
} -> session_persist_check
print "sessionStorage['session_persist_key'] after navigation:"
print session_persist_check

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ All storage tests completed!"