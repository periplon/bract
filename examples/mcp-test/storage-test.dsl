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
print "=== Testing Cookies ==="

# Set a cookie
call browser_set_cookie {
  tabId: tab.id,
  name: "test_cookie",
  value: "test_value_123",
  domain: "example.com"
} -> setCookieResult

assert setCookieResult.success == true, "Failed to set cookie"
print "✓ Cookie set successfully"

# Get cookies
call browser_get_cookies {
  tabId: tab.id
} -> cookies

# Find our test cookie
set found = false
loop cookie in cookies {
  if cookie.name == "test_cookie" {
    assert cookie.value == "test_value_123", "Cookie value mismatch"
    set found = true
  }
}
assert found == true, "Test cookie not found"
print "✓ Cookie retrieved successfully"

# Delete the cookie
call browser_delete_cookies {
  name: "test_cookie",
  domain: "example.com"
}
print "✓ Cookie deleted"

# Test localStorage
print "\n=== Testing localStorage ==="

call browser_set_local_storage {
  tabId: tab.id,
  key: "testKey",
  value: "testValue"
}
print "✓ localStorage item set"

call browser_get_local_storage {
  tabId: tab.id,
  key: "testKey"
} -> localValue

assert localValue == "testValue", "localStorage value mismatch"
print "✓ localStorage item retrieved"

# Test sessionStorage
print "\n=== Testing sessionStorage ==="

call browser_set_session_storage {
  tabId: tab.id,
  key: "sessionKey",
  value: "sessionValue"
}
print "✓ sessionStorage item set"

call browser_get_session_storage {
  tabId: tab.id,
  key: "sessionKey"
} -> sessionValue

assert sessionValue == "sessionValue", "sessionStorage value mismatch"
print "✓ sessionStorage item retrieved"

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ All storage tests passed!"
