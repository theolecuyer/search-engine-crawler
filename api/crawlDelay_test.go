package main

import (
	"testing"
	"time"
)

func TestCrawlDelay(t *testing.T) {
	indx := MakeInMemoryIndex()
	go webserver(indx)
	time.Sleep(1 * time.Second)
	tests := []struct {
		testName string
		url      string
		delay    int
	}{
		{
			"DelayTest1",
			"http://localhost:8080/top10/Dracula%20%7C%20Project%20Gutenberg/index.html",
			1,
		},
	}
	{
		for _, test := range tests {
			t1 := time.Now()
			idx := MakeInMemoryIndex()
			crawl(test.url, idx)
			t2 := time.Now()
			if t2.Sub(t1) < (2 * time.Second) {
				t.Errorf("Delay was too fast\n")
			}
		}
	}
}
