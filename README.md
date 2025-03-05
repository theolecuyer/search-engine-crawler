# search-engine-crawler

**search-engine-crawler** is a focused web crawler, indexer, and search engine designed to search and rank content within a single website. Using ethical crawling practices, it reads and respects `robots.txt` directives for allowed paths and crawl delay and calculates search relevance through **TF-IDF** scoring to deliver relevant results.

Currently the crawler has completely crawled and indexed my schools website:
**[https://www.usfca.edu/](https://www.usfca.edu/)**
This is all of the current data available to search which contains over:
**5500 links**
**41,709,334 possible relevant search terms**

## Features

- **Ethical Crawling**: Complies with `robots.txt` files to respect website rules.
- **TF-IDF Ranking**: Ranks pages using TF-IDF (Term Frequency-Inverse Document Frequency) to ensure relevant search results.
- **Full-Text Search**: Search through content from the website with ranked results.
- **One-Site Focused**: Indexes and searches only a single website, staying within the boundaries of the given host.

## Used Technologies

**Frameworks, Databases, and Languages:**

- **Next.js**: For handling the front-end, user interface, and calls to the custom API.
- **PostgreSQL**: Storage of crawled data and search indexes. Utilized for its concurrency power.
- **Go**: For building the custom API and handling backend indexing and searching processes.

**Key Libraries and Packages:**

- [**kljensen/snowball**](https://github.com/kljensen/snowball): Used for stemming in search queries, enhancing the accuracy of TF-IDF calculations.
- [**lib/pq**](https://github.com/lib/pq): Go client for PostgreSQL, allowing efficient data handling.
- [**Diggernaut/stopwords**](https://github.com/Diggernaut/stopwords): For removing common stopwords from text, improving search relevance.
- [**golang.org/x/net/html**](https://pkg.go.dev/golang.org/x/net/html): Provides HTML parsing for text extraction during crawling.
- [**google/uuid**](https://github.com/google/uuid): Generates unique IDs for indexed pages, ensuring reliable and unique identification within the database.

## How It Works

1. **Crawling**: Reads `robots.txt` to identify pages it is allowed to crawl. It then traverses the pages of the given website, extracting text content for indexing.
2. **Indexing**: Text content is stored in a PostgreSQL database, with keywords stemmed and stopwords removed for efficient and relevant indexing.
3. **Searching**: Using TF-IDF scoring, search results are ranked to show the most relevant pages based on the user's query.

## Running the project

Currently the front-end of project is still under development, but you can view the progress and try it for yourself at
https://sitesearch.vercel.app/

Last Updated: **03.05.2025**
