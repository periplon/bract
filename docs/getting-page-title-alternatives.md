# Alternative Methods to Get Page Title Without browser_execute_script

This document outlines alternative methods to get the page title without using `browser_execute_script`.

## Method 1: Using browser_list_tabs

The `browser_list_tabs` tool returns an array of tab objects, each containing a `title` field.

### Implementation Details

The `Tab` struct (defined in `/internal/browser/types.go`) includes:
```go
type Tab struct {
    ID      int    `json:"id"`
    URL     string `json:"url"`
    Title   string `json:"title"`  // <-- Page title is available here
    Active  bool   `json:"active"`
    Index   int    `json:"index"`
    Favicon string `json:"favicon,omitempty"`
}
```

### Usage Example

```dsl
# List all tabs to get their titles
call browser_list_tabs -> tabs

# Access title from the first tab
print "First tab title: " + tabs[0].title

# Find a specific tab by ID and get its title
loop tab in tabs {
    if tab.id == targetTabId {
        print "Tab title: " + tab.title
    }
}
```

### Advantages
- No JavaScript execution required
- Gets titles for all tabs in one call
- Includes other useful metadata (URL, active state, etc.)
- More efficient for getting multiple tab titles

### Limitations
- Returns data for all tabs, not just one
- Requires filtering if you only need one specific tab's title

## Method 2: Check browser_navigate Response

The `browser_navigate` tool might include page information in its response, though this needs verification based on the browser extension's implementation.

### Current Implementation
Based on the code analysis, `browser_navigate` returns a raw JSON response from the browser extension:

```go
// Navigate navigates to a URL in a tab
func (c *Client) Navigate(ctx context.Context, tabID int, url string, waitUntilLoad bool) (json.RawMessage, error) {
    // ... implementation
    response, err := c.sendCommand(ctx, "navigate", params)
    return response, err
}
```

The actual response structure depends on what the Chrome extension returns, which typically includes:
```json
{
    "success": true,
    // Additional fields may be included by the extension
}
```

## Method 3: Using browser_extract_content for Title Element

You can extract the title using the `<title>` tag selector:

```dsl
# Extract the title element content
call browser_extract_content {
    tabId: tab.id,
    selector: "title",
    contentType: "text"
} -> titleContent

print "Page title: " + titleContent[0]
```

### Advantages
- Works with standard HTML selectors
- No JavaScript execution needed
- Can be combined with other content extraction

### Limitations
- Returns an array (even for single elements)
- Might not work if the title is dynamically set

## Recommended Approach

**For most use cases, use `browser_list_tabs`:**

1. It's the most reliable method
2. Provides additional metadata
3. No JavaScript execution required
4. Works consistently across all pages

### Example Implementation

```dsl
# Function to get title for a specific tab
define getTabTitle(tabId) {
    call browser_list_tabs -> tabs
    
    loop tab in tabs {
        if tab.id == tabId {
            return tab.title
        }
    }
    
    return "Tab not found"
}

# Usage
set title = getTabTitle(myTab.id)
print "Current page title: " + title
```

## Performance Comparison

| Method | Performance | Reliability | Use Case |
|--------|------------|-------------|----------|
| `browser_list_tabs` | Fast | High | Best for getting titles of one or more tabs |
| `browser_extract_content` | Medium | Medium | When you need other page content too |
| `browser_execute_script` | Slow | High | When you need dynamic title or complex logic |

## Conclusion

The `browser_list_tabs` method is the recommended approach for getting page titles as it:
- Doesn't require JavaScript execution
- Is more performant
- Provides additional useful metadata
- Is more reliable across different page states