# Storage Example
# Demonstrates practical use of cookies, localStorage, and sessionStorage

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

# Create a tab and navigate to a site
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://example.com"
}

print "=== Storage Example ==="

# Example 1: Using cookies for authentication
print "\n1. Cookie Authentication Example:"
call browser_set_cookie {
  name: "auth_token",
  value: "abc123xyz",
  domain: ".example.com",
  path: "/",
  secure: true,
  httpOnly: true,
  expirationDate: 1735689600  # Jan 1, 2025
}
print "✓ Set authentication cookie"

# Example 2: Using localStorage for user preferences
print "\n2. User Preferences Example:"
call browser_set_local_storage {
  tabId: tab.id,
  key: "theme",
  value: "dark"
}
call browser_set_local_storage {
  tabId: tab.id,
  key: "language",
  value: "en"
}
print "✓ Saved user preferences to localStorage"

# Example 3: Using sessionStorage for temporary data
print "\n3. Temporary Data Example:"
call browser_set_session_storage {
  tabId: tab.id,
  key: "form_draft",
  value: "{\"title\":\"My Draft\",\"content\":\"Work in progress...\"}"
}
print "✓ Saved form draft to sessionStorage"

# Retrieve and display all storage
print "\n=== Current Storage State ==="

# Get cookies
call browser_get_cookies {url: "https://example.com"} -> cookies
print "Cookies:"
print cookies

# Get localStorage preferences
call browser_get_local_storage {tabId: tab.id, key: "theme"} -> theme
call browser_get_local_storage {tabId: tab.id, key: "language"} -> lang
print "User preferences: theme=" + theme.value + ", language=" + lang.value

# Get sessionStorage draft
call browser_get_session_storage {tabId: tab.id, key: "form_draft"} -> draft
print "Form draft:"
print draft.value

# Clean up specific items
print "\n=== Cleanup Example ==="
call browser_delete_cookies {name: "auth_token"}
print "✓ Logged out (deleted auth cookie)"

call browser_clear_session_storage {tabId: tab.id}
print "✓ Cleared temporary session data"

# Close tab
call browser_close_tab {tabId: tab.id}
print "\n✓ Example completed!"