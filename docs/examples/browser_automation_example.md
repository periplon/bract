# Browser Automation Examples

This document provides examples of using the MCP Browser Automation Server.

## Basic Navigation

### Navigate to a website
```
Use browser_navigate to go to https://www.example.com
```

### Take a screenshot
```
Take a screenshot of the current page using browser_screenshot
```

## Form Interaction

### Fill and submit a form
```
1. Navigate to the login page
2. Type username in the #username field
3. Type password in the #password field
4. Click the submit button
```

## Tab Management

### Work with multiple tabs
```
1. List all open tabs
2. Create a new tab with https://google.com
3. Switch back to the first tab
4. Close the second tab
```

## Content Extraction

### Extract text from a page
```
1. Navigate to a news website
2. Extract all headlines (h1, h2 elements)
3. Extract the main article text
```

## Advanced Automation

### Wait for dynamic content
```
1. Navigate to a page with lazy-loaded content
2. Wait for the element .lazy-content to appear
3. Scroll to the bottom of the page
4. Extract the newly loaded content
```

### Execute custom JavaScript
```
1. Navigate to a page
2. Execute JavaScript to get the page title
3. Execute JavaScript to count all links on the page
```

## Storage Operations

### Work with cookies
```
1. Get all cookies for the current domain
2. Set a new cookie named "session_id"
3. Delete specific cookies
```

### Local storage manipulation
```
1. Set a value in localStorage
2. Read the value back
3. Clear specific localStorage items
```

## Complete Workflow Example

### Automated form submission with verification
```
1. Navigate to https://example-form.com
2. Wait for the form to load (element #contact-form)
3. Type "John Doe" in the input[name="fullname"] field
4. Type "john@example.com" in the input[name="email"] field
5. Type "This is a test message" in the textarea[name="message"] field
6. Take a screenshot before submission
7. Click the button[type="submit"]
8. Wait for the success message (element .success-message)
9. Extract the confirmation text
10. Take a final screenshot
```

## Error Handling Examples

### Handle navigation failures
```
1. Try to navigate to an invalid URL
2. Handle the error gracefully
3. Navigate to a fallback URL
```

### Handle missing elements
```
1. Try to click on a non-existent element
2. Wait for an element with timeout
3. Proceed with alternative action if element not found
```

## Performance Testing

### Measure page load time
```
1. Record the start time
2. Navigate to a website
3. Wait for the page to fully load
4. Execute JavaScript to get performance metrics
5. Calculate and report load time
```

## Integration Patterns

### Data collection workflow
```
1. Read a list of URLs from input
2. For each URL:
   - Navigate to the URL
   - Wait for content to load
   - Extract specific data
   - Store results
3. Compile and return all results
```

### Authentication flow
```
1. Navigate to login page
2. Check if already logged in (look for logout button)
3. If not logged in:
   - Fill in credentials
   - Submit form
   - Wait for redirect
4. Verify successful login
5. Proceed with authenticated actions
```