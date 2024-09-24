package main

import "time"

// Type used across all files
type Frequency map[string]int   //Maps links to their word frequency
type Index map[string]Frequency //Maps each word and their correpsonding links and frequencies

func main() {
	dataWordFreq := make(Index)
	dataDocLen := make(Frequency)
	go webserver(dataWordFreq, dataDocLen)
	//index("https://cs272-f24.github.io/tests/rnj/", dataWordFreq, dataDocLen)
	index("http://localhost:8080/top10/index.html", dataWordFreq, dataDocLen)
	for {
		//Run the server until manual exit
		time.Sleep(1 * time.Second)
	}
}
