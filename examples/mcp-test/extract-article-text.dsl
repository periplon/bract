# Extract and analyze article text from news/blog websites
# This example shows how to extract clean text from articles
# and perform basic text analysis

# Connect to browser extension server
connect "./bin/mcp-browser-server"

# Wait for browser extension to connect
call browser_wait_for_connection {
    timeout: 5
} -> connection_result
print "Browser connection status:"
print connection_result

# Common article selectors for popular websites
set articleSelectors = "article, main article, .post-content, .entry-content, .article-body, .story-body, [itemprop='articleBody']"

# Extract article title
print "\nExtracting article title..."
call browser_extract_text {
    selector: "h1, article h1, .article-title, .post-title, [itemprop='headline']"
} -> title
if len(title) > 0 {
    print "Title found: " + str(len(title)) + " characters"
} else {
    print "No title found"
}

# Extract article content
print "\nExtracting article content..."
call browser_extract_text {
    selector: articleSelectors
} -> articleText

if len(articleText) > 0 {
    # Basic text analysis
    print "\nArticle Analysis:"
    print "- Character count: " + str(len(articleText))
    
    # Simple word count estimation (rough approximation)
    # Assuming average word length of 5 characters + 1 space
    set estimatedWords = len(articleText) / 6
    print "- Estimated word count: " + str(estimatedWords)
    print "- Estimated reading time: " + str(estimatedWords / 200) + " minutes"
    
    print "\nArticle text extracted successfully!"
} else {
    print "No article content found. This might not be an article page."
    print "Trying to extract all text content..."
    
    # Fallback: extract all text
    call browser_extract_text -> allText
    print "Total text extracted: " + str(len(allText)) + " characters"
}

# Extract metadata if available
print "\nExtracting metadata..."

# Extract author
call browser_extract_text {
    selector: ".author, .by-author, .article-author, [itemprop='author'], .author-name"
} -> author
if len(author) > 0 {
    print "Author found: " + str(len(author)) + " characters"
} else {
    print "No author information found"
}

# Extract date
call browser_extract_text {
    selector: "time, .publish-date, .article-date, [itemprop='datePublished'], .date"
} -> date
if len(date) > 0 {
    print "Date found: " + str(len(date)) + " characters"
} else {
    print "No date information found"
}

# Extract categories or tags
call browser_extract_text {
    selector: ".category, .tag, .tags, [rel='tag'], .article-tags"
} -> tags
if len(tags) > 0 {
    print "Tags/Categories found: " + str(len(tags)) + " characters"
}

print "\n✓ Article extraction complete!"