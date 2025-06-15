// Extract and analyze article text from news/blog websites
// This example shows how to extract clean text from articles
// and perform basic text analysis

// Connect to browser extension
connect

// Common article selectors for popular websites
articleSelectors = "article, main article, .post-content, .entry-content, .article-body, .story-body, [itemprop='articleBody']"

// Extract article title
print("Extracting article title...")
title -> browser_extract_text({
    selector: "h1, article h1, .article-title, .post-title, [itemprop='headline']"
})
print("Title: " + title)

// Extract article content
print("\nExtracting article content...")
articleText -> browser_extract_text({
    selector: articleSelectors
})

if (len(articleText) > 0) {
    // Basic text analysis
    wordCount = len(articleText.split(" "))
    charCount = len(articleText)
    
    print("Article Analysis:")
    print("- Word count: " + str(wordCount))
    print("- Character count: " + str(charCount))
    print("- Estimated reading time: " + str(wordCount / 200) + " minutes")
    
    // Extract first paragraph as summary
    paragraphs = articleText.split("\n\n")
    if (len(paragraphs) > 0) {
        print("\nFirst paragraph (summary):")
        print(paragraphs[0])
    }
    
    // Save to variable for further processing
    print("\nArticle text extracted successfully!")
    print("Text is now available in 'articleText' variable")
    
    // Example: Search for keywords
    keywords = ["AI", "technology", "innovation", "future"]
    print("\nKeyword analysis:")
    for (keyword in keywords) {
        if (articleText.indexOf(keyword) != -1) {
            print("- Found keyword: " + keyword)
        }
    }
} else {
    print("No article content found. This might not be an article page.")
    print("Trying to extract all text content...")
    
    // Fallback: extract all text
    allText -> browser_extract_text()
    print("Total text extracted: " + str(len(allText)) + " characters")
}

// Extract metadata if available
print("\nExtracting metadata...")
author -> browser_extract_text({
    selector: ".author, .by-author, .article-author, [itemprop='author']"
})
if (len(author) > 0) {
    print("Author: " + author)
}

date -> browser_extract_text({
    selector: "time, .publish-date, .article-date, [itemprop='datePublished']"
})
if (len(date) > 0) {
    print("Date: " + date)
}

print("\nExtraction complete!")