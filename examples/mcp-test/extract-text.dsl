// Extract plain text from current browser tab
// This example demonstrates using the browser_extract_text tool
// to get clean, readable text content from web pages

// Connect to browser extension
connect

// Extract text from the entire page (body element)
print("Extracting text from entire page...")
pageText -> browser_extract_text()
print("Page text length: " + str(len(pageText)))
print("First 500 characters:")
print(pageText[0:500] + "...")

// Extract text from specific elements
print("\n\nExtracting text from specific elements...")

// Extract main content (adjust selector based on the website)
mainContent -> browser_extract_text({
    selector: "main, article, [role='main'], .content, #content"
})
if (len(mainContent) > 0) {
    print("Main content found (" + str(len(mainContent)) + " characters)")
    print("Preview: " + mainContent[0:200] + "...")
} else {
    print("No main content found with common selectors")
}

// Extract all headings
print("\n\nExtracting all headings...")
headings -> browser_extract_text({
    selector: "h1, h2, h3, h4, h5, h6"
})
print("Headings found:")
print(headings)

// Extract navigation links
print("\n\nExtracting navigation text...")
navText -> browser_extract_text({
    selector: "nav, header, .navigation, .nav, .menu"
})
if (len(navText) > 0) {
    print("Navigation text (" + str(len(navText)) + " characters):")
    print(navText[0:200] + "...")
}

// Extract specific content by ID or class
// Uncomment and adjust selector for your specific use case
// specificContent -> browser_extract_text({
//     selector: "#specific-id, .specific-class"
// })
// print("Specific content: " + specificContent)

print("\n\nText extraction complete!")