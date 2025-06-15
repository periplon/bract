// Compare browser_extract_content vs browser_extract_text
// This example shows the difference between extracting HTML content
// and plain text from web pages

// Connect to browser extension
connect

print("Comparing content extraction methods...")
print("=====================================\n")

// Define selector for comparison
selector = "p"  // Extract all paragraphs

// Method 1: Extract HTML content
print("1. Using browser_extract_content (HTML):")
htmlContent -> browser_extract_content({
    selector: selector,
    type: "html"
})
print("HTML content (first item):")
if (len(htmlContent) > 0) {
    print(htmlContent[0])
    print("Total items: " + str(len(htmlContent)))
}

// Method 2: Extract as plain text using browser_extract_content
print("\n\n2. Using browser_extract_content (text mode):")
textArray -> browser_extract_content({
    selector: selector,
    type: "text"
})
print("Text array (first 3 items):")
for (i in [0, 1, 2]) {
    if (i < len(textArray)) {
        print("  [" + str(i) + "]: " + textArray[i])
    }
}
print("Total items: " + str(len(textArray)))

// Method 3: Extract as plain text using browser_extract_text
print("\n\n3. Using browser_extract_text (new tool):")
plainText -> browser_extract_text({
    selector: selector
})
print("Plain text (first 500 chars):")
print(plainText[0:500] + "...")
print("Total length: " + str(len(plainText)) + " characters")

// Key differences
print("\n\nKey Differences:")
print("================")
print("- browser_extract_content returns an array of strings (one per element)")
print("- browser_extract_text returns a single string with all text combined")
print("- browser_extract_text strips ALL HTML tags and cleans up formatting")
print("- browser_extract_text is ideal for reading content, text analysis, etc.")

// Practical example: Word count
print("\n\nPractical Example - Word Count:")
words = plainText.split(" ")
wordCount = len(words)
print("Total words in all paragraphs: " + str(wordCount))

// Find longest word
longestWord = ""
for (word in words) {
    if (len(word) > len(longestWord)) {
        longestWord = word
    }
}
print("Longest word: " + longestWord + " (" + str(len(longestWord)) + " chars)")

print("\n\nComparison complete!")