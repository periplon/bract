# Cookie Management Test
# Comprehensive demonstration of cookie operations including deletion

# Connect to the MCP browser server
connect "./bin/mcp-browser-server"

# Wait for browser extension
call browser_wait_for_connection {timeout: 5}

# Create a tab and navigate
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://example.com"
}

print "=== Cookie Management Test ==="

print "\n1. Setting multiple cookies:"
# Set various types of cookies
call browser_set_cookie {
  name: "session_id",
  value: "abc123xyz",
  domain: ".example.com",
  path: "/",
  secure: true,
  httpOnly: true
} -> cookie1
print "✓ Set session cookie"

call browser_set_cookie {
  name: "user_pref",
  value: "theme=dark;lang=en",
  domain: ".example.com",
  path: "/",
  secure: false,
  httpOnly: false,
  expirationDate: 1893456000  # Jan 1, 2030
} -> cookie2
print "✓ Set preference cookie with expiration"

call browser_set_cookie {
  name: "analytics_id",
  value: "GA123456",
  domain: "example.com",
  path: "/",
  secure: true,
  httpOnly: false
} -> cookie3
print "✓ Set analytics cookie"

call browser_set_cookie {
  name: "temp_data",
  value: "temporary",
  domain: ".example.com",
  path: "/temp",
  secure: false,
  httpOnly: false
} -> cookie4
print "✓ Set path-specific cookie"

print "\n2. Retrieving all cookies:"
# Get all cookies for the domain
call browser_get_cookies {
  url: "https://example.com"
} -> allCookies
print "Total cookies for example.com: " + str(len(allCookies))
loop cookie in allCookies {
  print "  - " + cookie.name + " = " + cookie.value
}

print "\n3. Testing specific cookie deletion:"
# Delete a specific cookie by name (with proper path)
call browser_delete_cookies {
  url: "https://example.com/temp",
  name: "temp_data"
} -> deleteResult1
print "✓ Deleted 'temp_data' cookie"

# Verify deletion (check all cookies with the correct path)
call browser_get_cookies {
  url: "https://example.com/temp"
} -> checkDeleted
set found = false
loop c in checkDeleted {
  if c.name == "temp_data" {
    set found = true
  }
}
assert found == false, "Cookie should be deleted"
print "✓ Confirmed 'temp_data' cookie is deleted"

print "\n4. Testing deletion of non-existent cookie:"
# Try to delete a cookie that doesn't exist
call browser_delete_cookies {
  url: "https://example.com",
  name: "non_existent_cookie"
} -> deleteResult2
print "✓ Deletion of non-existent cookie handled gracefully"

print "\n5. Setting and deleting httpOnly cookie:"
# Set an httpOnly cookie
call browser_set_cookie {
  name: "secure_token",
  value: "secret123",
  domain: ".example.com",
  path: "/",
  secure: true,
  httpOnly: true
} -> secureCookie
print "✓ Set httpOnly cookie"

# Delete the httpOnly cookie
call browser_delete_cookies {
  url: "https://example.com",
  name: "secure_token"
} -> deleteSecure
print "✓ Deleted httpOnly cookie"

# Verify it's gone
call browser_get_cookies {
  url: "https://example.com"
} -> checkSecure
set secureFound = false
loop c in checkSecure {
  if c.name == "secure_token" {
    set secureFound = true
  }
}
assert secureFound == false, "Secure cookie should be deleted"
print "✓ Confirmed httpOnly cookie is deleted"

print "\n6. Testing URL-based cookie deletion:"
# Navigate to a different subdomain
call browser_navigate {
  tabId: tab.id,
  url: "https://subdomain.example.com"
}

# Set a cookie specific to subdomain
call browser_set_cookie {
  name: "subdomain_cookie",
  value: "sub_value",
  domain: "subdomain.example.com",
  path: "/"
} -> subCookie
print "✓ Set subdomain-specific cookie"

# Delete cookies for specific URL
call browser_delete_cookies {
  url: "https://subdomain.example.com",
  name: "subdomain_cookie"
} -> deleteSubResult
print "✓ Deleted subdomain cookie"

print "\n7. Testing bulk cookie operations:"
# Get remaining cookies
call browser_get_cookies {
  url: "https://example.com"
} -> remainingCookies
print "Remaining cookies: " + str(len(remainingCookies))

# Delete cookies one by one
loop cookie in remainingCookies {
  call browser_delete_cookies {
    url: "https://example.com",
    name: cookie.name
  } -> deleteLoop
  print "✓ Deleted cookie: " + cookie.name
}

# Verify all cookies are deleted
call browser_get_cookies {
  url: "https://example.com"
} -> finalCheck
print "\nFinal cookie count: " + str(len(finalCheck))
# Work around DSL type comparison issue - check if array is empty
set allDeleted = true
loop c in finalCheck {
  set allDeleted = false
}
assert allDeleted == true, "All cookies should be deleted"
print "✓ All cookies successfully deleted"

print "\n8. Testing cookie deletion without URL:"
# Set a new cookie
call browser_set_cookie {
  name: "test_final",
  value: "final_value",
  domain: ".example.com",
  path: "/"
}

# Delete without specifying URL (should use current page URL)
call browser_delete_cookies {
  name: "test_final"
} -> deleteFinal
print "✓ Deleted cookie using current page URL"

# Clean up
call browser_close_tab {tabId: tab.id}

print "\n✓ Cookie management test completed!"