# MCP Test DSL Documentation

The MCP Test DSL is a domain-specific language designed for testing Model Context Protocol (MCP) servers. It provides a simple, readable syntax for writing test scripts and reusable automations.

## Overview

The DSL allows you to:
- Connect to MCP servers
- Call tools and verify responses
- Write assertions and test logic
- Create reusable automation scripts
- Test browser automation capabilities

## Installation

Build the test runner:
```bash
go build -o mcp-test ./cmd/mcp-test
```

## Usage

Run a DSL script:
```bash
./mcp-test script.dsl
```

Validate syntax without executing:
```bash
./mcp-test -validate script.dsl
```

Format a script:
```bash
./mcp-test -format script.dsl
```

## Language Reference

### Comments
```dsl
# This is a comment
```

### Variables and Types
```dsl
# Set variables
set name = "value"
set number = 42
set decimal = 3.14
set boolean = true
set array = [1, 2, 3]
set object = {key: "value", count: 10}

# Access fields
set value = object.key
set item = array[0]
```

### Connecting to MCP Servers
```dsl
# Basic connection
connect "/path/to/mcp-server"

# With arguments
connect "./server" "--debug" "--port" "8080"

# With options
connect "./server" {
  timeout: 30
  retries: 3
}
```

### Calling Tools
```dsl
# Simple tool call
call tool_name

# With arguments
call navigate {url: "https://example.com"}

# Store result in variable
call list_tools -> tools

# With complex arguments
call execute_script {
  tabId: tab.id,
  script: "return document.title"
} -> result
```

### Control Flow

#### If Statements
```dsl
if condition {
  # then block
} else if other_condition {
  # else if block
} else {
  # else block
}
```

#### Loops
```dsl
# Loop over array
loop item in array {
  print item
}

# Loop over object (gives key-value pairs)
loop entry in object {
  print entry.key + ": " + entry.value
}
```

### Assertions
```dsl
# Basic assertion
assert condition

# With custom message
assert result.success == true, "Operation should succeed"
```

### Waiting
```dsl
# Wait for condition
wait element.visible == true

# With timeout (seconds)
wait page.loaded == true, 10

# With timeout and interval (milliseconds)
wait counter > 100, 30, 500
```

### Functions
```dsl
# Built-in functions
len(array)        # Length of array/string/object
str(value)        # Convert to string
int(value)        # Convert to integer
float(value)      # Convert to float
json(value)       # Convert to JSON string
```

### Operators
```dsl
# Arithmetic
+ - * /

# Comparison
== != < > <= >=

# Logical
&& || !

# String concatenation
"Hello " + "World"
```

### Defining Automations
```dsl
# Define reusable automation
define login(username, password) {
  call navigate {url: "/login"}
  call type {selector: "#username", text: username}
  call type {selector: "#password", text: password}
  call click {selector: "#submit"}
}

# Run automation
run login("user@example.com", "secret123")
```

### Printing Output
```dsl
print "Hello, World!"
print "Value: " + str(value)
print result
```

## Example Scripts

### Basic Connectivity Test
```dsl
# Connect to server
connect "./mcp-server"

# List available tools
call list_tools -> tools
print "Available tools: " + str(len(tools))

# Verify tools exist
assert len(tools) > 0, "Server should have tools"
```

### Browser Automation Test
```dsl
# Connect to browser server
connect "./mcp-browser-server"

# Create and navigate tab
call create_tab -> tab
call navigate {
  tabId: tab.id,
  url: "https://example.com"
}

# Wait for page load
wait 2

# Extract title
call execute_script {
  tabId: tab.id,
  script: "document.title"
} -> title

print "Page title: " + title.result

# Clean up
call close_tab {tabId: tab.id}
```

### Reusable Test Pattern
```dsl
# Define page object pattern
define TestPage(url) {
  call create_tab -> tab
  call navigate {tabId: tab.id, url: url}
  
  define click(selector) {
    call click {tabId: tab.id, selector: selector}
  }
  
  define type(selector, text) {
    call type {
      tabId: tab.id,
      selector: selector,
      text: text
    }
  }
  
  define cleanup() {
    call close_tab {tabId: tab.id}
  }
}

# Use the pattern
run TestPage("https://example.com")
run type("#search", "test query")
run click("#submit")
wait 2
run cleanup()
```

## Testing Best Practices

1. **Use descriptive variable names**
   ```dsl
   call create_tab -> mainTab  # Good
   call create_tab -> t        # Less clear
   ```

2. **Add assertions for critical checks**
   ```dsl
   call navigate {url: url} -> result
   assert result.success == true, "Navigation must succeed"
   ```

3. **Clean up resources**
   ```dsl
   call create_tab -> tab
   # ... do work ...
   call close_tab {tabId: tab.id}  # Always clean up
   ```

4. **Use automations for repeated tasks**
   ```dsl
   define setup_test_tab(url) {
     call create_tab -> tab
     call navigate {tabId: tab.id, url: url}
     set current_tab = tab
   }
   ```

5. **Handle errors gracefully**
   ```dsl
   call risky_operation -> result
   if result.error {
     print "Operation failed: " + result.error
     # Handle error...
   }
   ```

## Integration with Go Tests

You can use the MCP client directly in Go tests:

```go
import (
    "testing"
    "github.com/periplon/bract/internal/mcpclient"
)

func TestMCPServer(t *testing.T) {
    harness := mcpclient.NewTestHarness(t, "./mcp-server")
    harness.RunTest(func(client *mcpclient.TestClient) {
        // List tools
        tools := client.MustListTools(context.Background())
        assert.NotEmpty(t, tools)
        
        // Call a tool
        result := client.MustCallTool(ctx, "echo", map[string]string{
            "message": "hello",
        })
        
        // Make assertions
        client.AssertToolResult(ctx, "echo", args, expected)
    })
}
```

## Troubleshooting

### Script Validation Errors
- Check syntax with `-validate` flag
- Ensure all braces and quotes are balanced
- Verify variable names don't use reserved keywords

### Connection Issues
- Ensure MCP server is executable
- Check server path is correct
- Verify server implements MCP protocol correctly

### Tool Call Failures
- List available tools first
- Check tool arguments match expected schema
- Verify tool permissions and capabilities