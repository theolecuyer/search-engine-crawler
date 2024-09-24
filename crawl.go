package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/kljensen/snowball"
)

func crawl(baseURL string) map[string][]string {
	//Format: "link": {"array", "of", "words"}
	allLinksWords := map[string][]string{}
	//Make a map for all visited urls
	visitedUrls := make(map[string]bool)
	host, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("URL Parse returned %v", err)
	}
	visitedUrls[baseURL] = true
	hostName := host.Host
	//Make a fifo "queue" of urls and initailize the "starter" url to it
	queue := []string{baseURL}

	for len(queue) > 0 {
		var currentUrl = queue[0] //Get the top element
		fmt.Println("Crawling at: " + currentUrl)
		queue = queue[1:] //"Pop" the top element
		if result, err := download(currentUrl); err != nil {
			log.Printf("Download returned: %v\n", err)
		} else {
			words, hrefs := extract(string(result))
			for i := range words {
				//Stem all words before inserting into map
				if stemmed, err := snowball.Stem(words[i], "english", true); err != nil {
					log.Printf("Snowball returned %v", err)
				} else {
					words[i] = stemmed
				}
			}
			allLinksWords[currentUrl] = words
			links := clean(baseURL, hrefs)
			for _, cleanedURL := range links {
				if !visitedUrls[cleanedURL.String()] && hostName == cleanedURL.Host {
					queue = append(queue, cleanedURL.String())
					visitedUrls[cleanedURL.String()] = true
				}
			}
		}
	}
	return allLinksWords
}
