# Extract plain text from current browser tab
# This example demonstrates using the browser_extract_text tool
# to get clean, readable text content from web pages

# Connect to browser extension server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
    timeout: 5
} -> connection_result
print "Browser connection status:"
print connection_result

# Extract text from the entire page (body element)
print "\nExtracting text from entire page..."
call browser_extract_text -> pageText
print "Page text length: " + str(len(pageText))

# Extract text from specific elements
print "\n\nExtracting text from specific elements..."

# Extract main content (adjust selector based on the website)
call browser_extract_text {
    selector: "main, article, [role='main'], .content, #content"
} -> mainContent

if len(mainContent) > 0 {
    print "Main content found (" + str(len(mainContent)) + " characters)"
} else {
    print "No main content found with common selectors"
}

# Extract all headings
print "\n\nExtracting all headings..."
call browser_extract_text {
    selector: "h1, h2, h3, h4, h5, h6"
} -> headings
print "Headings text length: " + str(len(headings))

# Extract navigation links
print "\n\nExtracting navigation text..."
call browser_extract_text {
    selector: "nav, header, .navigation, .nav, .menu"
} -> navText

if len(navText) > 0 {
    print "Navigation text found (" + str(len(navText)) + " characters)"
} else {
    print "No navigation text found"
}

# Extract paragraphs
print "\n\nExtracting paragraph text..."
call browser_extract_text {
    selector: "p"
} -> paragraphText
print "Paragraph text length: " + str(len(paragraphText))

print "\n\n✓ Text extraction complete!"