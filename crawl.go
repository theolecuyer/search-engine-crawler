package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/kljensen/snowball"
)

func crawl(baseURL string, index Indexes) {
	allLinksWords := make(map[string][]string) //Return: URL, slice of words
	visitedUrls := make(map[string]bool)       //Make a map for all visited urls
	stopWordMap := loadStopWords("stopwords-en.json")
	host, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("URL Parse returned %v", err)
	}
	visitedUrls[baseURL] = true
	hostName := host.Host
	queue := []string{baseURL} //FIFO queue

	for len(queue) > 0 {
		var currentUrl = queue[0] //Get the top element
		fmt.Println("Crawling at: " + currentUrl)
		queue = queue[1:] //"Pop" the top element
		if result, err := download(currentUrl); err != nil {
			log.Printf("Download returned: %v\n", err)
		} else {
			words, hrefs := extract(string(result))
			for _, word := range words {
				if stemmedWord, err := snowball.Stem(word, "english", true); err != nil {
					log.Printf("Snowball error: %v", err)
				} else {
					if _, exists := stopWordMap[stemmedWord]; !exists {
						allLinksWords[currentUrl] = append(allLinksWords[currentUrl], stemmedWord)
					}
				}
			}

			links := clean(baseURL, hrefs)
			for _, cleanedURL := range links {
				if !visitedUrls[cleanedURL.String()] && hostName == cleanedURL.Host {
					queue = append(queue, cleanedURL.String())
					visitedUrls[cleanedURL.String()] = true
				}
			}
		}
	}
	index.AddToIndex(allLinksWords)
}

func loadStopWords(link string) map[string]struct{} {
	stopWordsFile, err := os.Open(link)
	if err != nil {
		fmt.Println(err)
	}
	defer stopWordsFile.Close()
	byteValue, _ := io.ReadAll(stopWordsFile)
	var stopWords []string
	if json.Unmarshal(byteValue, &stopWords); err != nil {
		log.Printf("Json unmarshal returned: %v", err)
	}
	stopWordMap := make(map[string]struct{})
	for _, word := range stopWords {
		stopWordMap[word] = struct{}{}
	}
	return stopWordMap
}
