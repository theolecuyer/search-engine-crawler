package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kljensen/snowball"
)

func index(link string, dataWordFreq Index, dataDocLen Frequency) {
	//Load stopwords from the JSON file
	stopWordMap := loadStopWords("stopwords-en.json")
	words := crawl(link)
	for link, wordArray := range words {
		for _, word := range wordArray {
			if _, exists := stopWordMap[word]; !exists {
				if stemmedWord, err := snowball.Stem(word, "english", true); err != nil {
					log.Printf("Snowball error: %v", err)
				} else {
					if dataWordFreq[stemmedWord] == nil {
						dataWordFreq[stemmedWord] = make(Frequency)
					}
					dataWordFreq[stemmedWord][link]++
				}
			}
		}
		dataDocLen[link] = len(wordArray)
	}
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
