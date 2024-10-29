# SiteSearch

**SiteSearch** is a focused web crawler, indexer, and search tool designed to help you search and rank content within a single website. Using ethical crawling practices, SiteSearch reads and respects `robots.txt` directives for allowed paths and crawl delay and calculates search relevance through **TF-IDF** scoring to deliver relevant results.

## Features

- **Ethical Crawling**: Complies with `robots.txt` files to respect website rules.
- **TF-IDF Ranking**: Ranks pages using TF-IDF (Term Frequency-Inverse Document Frequency) to ensure relevant search results.
- **Full-Text Search**: Search through content from a given website with ranked results.
- **One-Site Focused**: Indexes and searches only a single website, making it an ideal tool for personal or single-site projects.

## Used Technologies

**Frameworks, Databases, and Languages:**

- **Next.js**: For handling the front-end, user interface, and calls to the custom API.
- **PostgreSQL**: Storage of crawled data and search indexes. Utilized for its concurrency power.
- **Go (Golang)**: For building the custom API and handling backend indexing and searching processes.

**Key Libraries and Packages:**

- [**kljensen/snowball**](https://github.com/kljensen/snowball): Used for stemming in search queries, enhancing the accuracy of TF-IDF calculations.
- [**lib/pq**](https://github.com/lib/pq): Go client for PostgreSQL, allowing efficient data handling.
- [**Diggernaut/stopwords**](https://github.com/Diggernaut/stopwords): For removing common stopwords from text, improving search relevance.
- [**golang.org/x/net/html**](https://pkg.go.dev/golang.org/x/net/html): Provides HTML parsing for text extraction during crawling.
- [**google/uuid**](https://github.com/google/uuid): Generates unique IDs for indexed pages, ensuring reliable and unique identification within the database.

## How It Works

1. **Crawling**: SiteSearch reads `robots.txt` to identify pages it is allowed to crawl. It then traverses the pages of the given website, extracting text content for indexing.
2. **Indexing**: Text content is stored in a PostgreSQL database, with keywords stemmed and stopwords removed for efficient and relevant indexing.
3. **Searching**: Using TF-IDF scoring, SiteSearch ranks search results to show the most relevant pages based on the user's query.

## Running SiteSearch

Currently the front-end of SiteSearch is still under development, but you can view the progress at
https://sitesearch.vercel.app/

Last Updated: **10.29.2024**
