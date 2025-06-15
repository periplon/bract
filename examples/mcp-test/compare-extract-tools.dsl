# Compare browser_extract_content vs browser_extract_text
# This example shows the difference between extracting HTML content
# and plain text from web pages

# Connect to browser extension server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
    timeout: 5
} -> connection_result
print "Browser connection status:"
print connection_result

print "\nComparing content extraction methods..."
print "====================================="

# Define selector for comparison
set selector = "p"  # Extract all paragraphs

# Method 1: Extract HTML content
print "\n1. Using browser_extract_content (HTML):"
call browser_extract_content {
    selector: selector,
    type: "html"
} -> htmlContent
print "HTML content items: " + str(len(htmlContent))
if len(htmlContent) > 0 {
    print "First item (HTML):"
    print htmlContent[0]
}

# Method 2: Extract as plain text using browser_extract_content
print "\n\n2. Using browser_extract_content (text mode):"
call browser_extract_content {
    selector: selector,
    type: "text"
} -> textArray
print "Text array items: " + str(len(textArray))
if len(textArray) > 0 {
    print "First 3 items:"
    set count = 0
    loop item in textArray {
        if count < 3 {
            print "  [" + str(count) + "]: " + item
            set count = count + 1
        }
    }
}

# Method 3: Extract as plain text using browser_extract_text
print "\n\n3. Using browser_extract_text (new tool):"
call browser_extract_text {
    selector: selector
} -> plainText
print "Plain text length: " + str(len(plainText)) + " characters"

# Key differences
print "\n\nKey Differences:"
print "================"
print "- browser_extract_content returns an array of strings (one per element)"
print "- browser_extract_text returns a single string with all text combined"
print "- browser_extract_text strips ALL HTML tags and cleans up formatting"
print "- browser_extract_text is ideal for reading content, text analysis, etc."

# Practical example: Character count comparison
print "\n\nPractical Example - Character Count:"
set totalCharsInArray = 0
loop item in textArray {
    set totalCharsInArray = totalCharsInArray + len(item)
}
print "Total chars in extract_content array: " + str(totalCharsInArray)
print "Total chars in extract_text: " + str(len(plainText))

# Show extraction from different element types
print "\n\nExtracting from different elements:"

# Links
call browser_extract_text {
    selector: "a"
} -> linkText
print "All link text combined: " + str(len(linkText)) + " characters"

# Lists
call browser_extract_text {
    selector: "li"
} -> listText
print "All list item text combined: " + str(len(listText)) + " characters"

print "\n\n✓ Comparison complete!"