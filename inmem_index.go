package main

import (
	"log"
	"sort"

	"github.com/kljensen/snowball"
)

type Indexes interface {
	AddToIndex(allWords map[string][]string)
	Search(query string) hits
}

type InMemoryIndex struct {
	wordFreq InvertedIndex
	doclen   Frequency
}

func MakeInMemoryIndex() *InMemoryIndex {
	return &InMemoryIndex{wordFreq: make(InvertedIndex), doclen: make(Frequency)}
}

func (i *InMemoryIndex) AddToIndex(allWords map[string][]string) {
	for url, words := range allWords {
		for _, word := range words {
			if i.wordFreq[word] == nil {
				i.wordFreq[word] = make(Frequency)
			}
			i.wordFreq[word][url]++
		}
		i.doclen[url] = len(words)
	}
}

func (i *InMemoryIndex) Search(query string) hits {
	results := hits{}
	resultFrequencies := Frequency{} // Return map "link" : # of word hits
	if stemmedWordQuery, err := snowball.Stem(query, "english", true); err == nil {
		if word, exists := i.wordFreq[stemmedWordQuery]; exists {
			for link, frequency := range word {
				tfIDFScore := tfIDF(frequency, i.doclen[link], len(i.doclen), len(i.wordFreq[link]))
				results = append(results, searchHit{link, frequency, tfIDFScore})
				resultFrequencies[link] = frequency
			}
		}
	} else {
		log.Printf("Snowball returned %v", err)
	}
	sort.Sort(results)
	return results
}
