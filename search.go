package main

import (
	"log"
	"math"
	"sort"

	"github.com/kljensen/snowball"
)

type searchHit struct {
	URL   string
	Count int
	tfIDF float64
}

type hits []searchHit

func (results hits) Len() int {
	return len(results)
}

func (results hits) Less(i, j int) bool {
	if results[i].tfIDF == results[j].tfIDF {
		return results[i].URL < results[j].URL
	}
	return results[i].tfIDF > results[j].tfIDF
}

func (results hits) Swap(i, j int) {
	results[i], results[j] = results[j], results[i]
}

func search(wordQuery string, dataWordFreq Index, dataDocLen Frequency) hits {
	results := hits{}
	resultFrequencies := Frequency{} // Return map "link" : # of word hits
	if stemmedWordQuery, err := snowball.Stem(wordQuery, "english", true); err == nil {
		if word, exists := dataWordFreq[stemmedWordQuery]; exists {
			for link, frequency := range word {
				tfIDFScore := tfIDF(link, frequency, dataDocLen, dataWordFreq, stemmedWordQuery)
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

func tfIDF(url string, wordFrequency int, dataDocLen Frequency, dataWordFreq Index, word string) float64 {
	tf := float64(wordFrequency) / float64(dataDocLen[url])
	idf := math.Log(float64(len(dataDocLen)) / float64(len(dataWordFreq[word])+1))
	tfIDFScore := (tf * idf)
	return tfIDFScore
}
