package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestDisallow(t *testing.T) {

	tests := []struct {
		testName   string
		url        string
		searchterm string
		expected   Frequency
	}{
		{
			"DisallowTest1",
			"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/index.html",
			"blood",
			map[string]int{
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap10.html": 15,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap24.html": 13,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap16.html": 8,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap03.html": 7,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap12.html": 7,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap11.html": 5,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap04.html": 5,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap13.html": 5,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap22.html": 4,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap15.html": 4,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap09.html": 3,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap14.html": 3,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap18.html": 3,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap02.html": 2,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap19.html": 2,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap25.html": 2,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap08.html": 2,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap27.html": 2,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap06.html": 1,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap20.html": 1,
				"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/chap07.html": 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			idx := MakeInMemoryIndex()
			Crawl(test.url, idx)
			got := idx.Search(test.searchterm)
			sort.Sort(got)
			results := Frequency{}
			for _, hit := range got {
				results[hit.URL] = hit.Count
			}
			if !reflect.DeepEqual(results, test.expected) {
				t.Errorf("Got: %v\nExpected: %v", results, test.expected)
			}
		})
	}
}
