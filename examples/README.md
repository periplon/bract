# Bract Examples

This directory contains comprehensive examples demonstrating the capabilities of the Bract MCP Browser Automation Server and DSL test runner.

## Overview

Bract provides two main components:
- **MCP Browser Server** (`mcp-browser-server`): A Model Context Protocol server that enables browser automation through a Chrome extension
- **MCP Test Runner** (`mcp-test`): A DSL interpreter for writing and running browser automation tests

## Prerequisites

1. Build the project binaries:
   ```bash
   make build
   ```

2. Install the Chrome extension (see main README for instructions)

3. Ensure the Chrome extension is running and connected

## Example Categories

### 1. Basic Examples

#### [basic.dsl](mcp-test/basic.dsl)
Tests basic server connectivity and tool listing.
```bash
./bin/mcp-test examples/mcp-test/basic.dsl
```

#### [basic-browser-wait.dsl](mcp-test/basic-browser-wait.dsl)
Demonstrates waiting for browser extension connection.
```bash
./bin/mcp-test examples/mcp-test/basic-browser-wait.dsl
```

### 2. Navigation Examples

#### [browser-navigation.dsl](mcp-test/browser-navigation.dsl)
Shows navigation, screenshots, and content extraction.
- Navigate to URLs
- Execute JavaScript
- Take screenshots
- Extract page content

#### [browser-reload.dsl](mcp-test/browser-reload.dsl)
Demonstrates page reload functionality.
- Normal reload (with cache)
- Hard reload (bypass cache)
- Storage persistence across reloads

### 3. Tab Management

#### [multi-tab.dsl](mcp-test/multi-tab.dsl)
Basic multi-tab operations.
- Create multiple tabs
- Navigate each tab
- List all tabs
- Close tabs

#### [tab-activation.dsl](mcp-test/tab-activation.dsl)
Advanced tab switching and activation.
- Switch between tabs
- Perform operations on non-active tabs
- Rapid tab switching
- Tab interaction after activation

### 4. Page Interaction

#### [form-interaction.dsl](mcp-test/form-interaction.dsl)
Form filling and submission automation.
- Fill input fields
- Click buttons
- Wait for elements
- Extract results

#### [browser-scroll.dsl](mcp-test/browser-scroll.dsl)
Comprehensive scrolling capabilities.
- Scroll to coordinates
- Smooth scrolling
- Scroll to elements
- Horizontal scrolling

### 5. Storage Management

#### [storage-test.dsl](mcp-test/storage-test.dsl)
Complete storage testing suite.
- Cookie operations
- localStorage management
- sessionStorage handling
- Storage persistence

#### [storage-example.dsl](mcp-test/storage-example.dsl)
Practical storage use cases.
- Authentication cookies
- User preferences
- Temporary data
- Cleanup operations

#### [cookie-management.dsl](mcp-test/cookie-management.dsl)
Advanced cookie operations.
- Set various cookie types
- Delete specific cookies
- Handle httpOnly cookies
- Bulk cookie operations

### 6. Advanced Automation

#### [advanced-automation.dsl](mcp-test/advanced-automation.dsl)
Complex automation patterns.
- Page Object pattern
- Reusable components
- Test suites
- Batch testing

## DSL Syntax Guide

### Basic Commands

```dsl
# Connect to server
connect "./bin/mcp-browser-server"

# Call a tool
call tool_name {
  param1: value1,
  param2: value2
} -> result

# Print output
print "Message: " + result.value

# Wait (in seconds)
wait 2

# Assertions
assert condition, "Error message"
```

### Variables and Data Types

```dsl
# Set variables
set myVar = "value"
set myNum = 42
set myBool = true
set myArray = [1, 2, 3]
set myObj = {key: "value"}

# Access properties
set tabId = tab.id
```

### Control Flow

```dsl
# Conditionals
if condition {
  # code
} else {
  # code
}

# Loops
loop item in array {
  print item
}

loop i in [1, 2, 3] {
  print "Index: " + str(i)
}
```

### Functions

```dsl
# Define reusable functions
define myFunction(param1, param2) {
  # function body
  set result = param1 + param2
}

# Call functions
run myFunction("arg1", "arg2")
```

## Available MCP Tools

### Tab Management
- `browser_list_tabs` - List all open tabs
- `browser_create_tab` - Create a new tab
- `browser_close_tab` - Close a tab
- `browser_activate_tab` - Switch to a tab

### Navigation
- `browser_navigate` - Navigate to a URL
- `browser_reload` - Reload the current page

### Interaction
- `browser_click` - Click on an element
- `browser_type` - Type text into a field
- `browser_scroll` - Scroll the page
- `browser_wait_for_element` - Wait for an element

### Content
- `browser_execute_script` - Execute JavaScript
- `browser_extract_content` - Extract page content
- `browser_screenshot` - Take a screenshot

### Storage
- `browser_get_cookies` - Get cookies
- `browser_set_cookie` - Set a cookie
- `browser_delete_cookies` - Delete cookies
- `browser_get_local_storage` - Get localStorage value
- `browser_set_local_storage` - Set localStorage value
- `browser_clear_local_storage` - Clear localStorage
- `browser_get_session_storage` - Get sessionStorage value
- `browser_set_session_storage` - Set sessionStorage value
- `browser_clear_session_storage` - Clear sessionStorage

### Connection
- `browser_wait_for_connection` - Wait for browser extension

## Running Examples

1. Start with the Chrome extension open
2. Run an example:
   ```bash
   ./bin/mcp-test examples/mcp-test/[example-name].dsl
   ```

## Writing Your Own Tests

1. Create a new `.dsl` file
2. Start with connecting to the server
3. Use the available tools to automate browser tasks
4. Add assertions to verify expected behavior
5. Clean up resources (close tabs) at the end

Example template:
```dsl
# My Test
# Description of what this test does

# Connect to server
connect "./bin/mcp-browser-server"

# Wait for browser
call browser_wait_for_connection {timeout: 5}

# Your test logic here
call browser_create_tab -> tab
call browser_navigate {
  tabId: tab.id,
  url: "https://example.com"
}

# Assertions
assert condition, "Error message"

# Cleanup
call browser_close_tab {tabId: tab.id}

print "✓ Test completed!"
```

## Troubleshooting

1. **Connection Issues**: Ensure the Chrome extension is installed and running
2. **Tool Not Found**: Check that you're using the correct tool name
3. **Timeout Errors**: Increase timeout values or add wait commands
4. **Element Not Found**: Verify CSS selectors are correct

## Contributing

When adding new examples:
1. Follow the existing naming convention
2. Include clear comments explaining the purpose
3. Add error handling where appropriate
4. Update this README with your example