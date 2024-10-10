package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	//Serve local files based off of the http request
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html-testfiles"+r.URL.Path)
	}))
	defer svr.Close()

	testCases := []struct {
		name     string
		url      string
		word     string
		expected Frequency
	}{
		{
			"searchVerona",
			svr.URL + "/tests/rnj/sceneI_30.0.html",
			"Verona",
			map[string]int{
				svr.URL + "/tests/rnj/sceneI_30.0.html": 1,
			},
		},
		{
			"searchBenvolio",
			svr.URL + "/tests/rnj/sceneI_30.1.html",
			"Benvolio",
			map[string]int{
				svr.URL + "/tests/rnj/sceneI_30.1.html": 26,
			},
		},
		{
			"searchRomeo",
			svr.URL + "/rnj.html",
			"Romeo",
			map[string]int{
				svr.URL + "/tests/rnj/sceneI_30.0.html":  2,
				svr.URL + "/tests/rnj/sceneI_30.1.html":  22,
				svr.URL + "/tests/rnj/sceneI_30.3.html":  2,
				svr.URL + "/tests/rnj/sceneI_30.4.html":  17,
				svr.URL + "/tests/rnj/sceneI_30.5.html":  15,
				svr.URL + "/tests/rnj/sceneII_30.2.html": 42,
				svr.URL + "/rnj.html":                    200,
				svr.URL + "/tests/rnj/sceneI_30.2.html":  15,
				svr.URL + "/tests/rnj/sceneII_30.0.html": 3,
				svr.URL + "/tests/rnj/sceneII_30.1.html": 10,
				svr.URL + "/tests/rnj/sceneII_30.3.html": 13,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			idx := MakeInMemoryIndex()
			crawl(test.url, idx)
			got := idx.Search(test.word)
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
