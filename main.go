package main

import (
	"database/sql"
	"flag"
	"log"
	"time"
)

// Type used across all files
type Frequency map[string]int           //Maps links to their word frequency
type InvertedIndex map[string]Frequency //Maps each word and their correpsonding links and frequencies

func main() {
	indexType := flag.String("index", "", "Specify the index type")
	flag.Parse()

	var indx Indexes
	if *indexType == "inmem" {
		indx = MakeInMemoryIndex()
		go webserver(indx)
	} else {
		db, err := sql.Open("sqlite", "index.db")
		if err != nil {
			log.Printf("DB open returned %v\n", err)
		}
		defer db.Close()
		indx = MakeDBIndex(db)
		go webserver(indx)
	}

	crawl("http://localhost:8080/top10/index.html", indx)
	//crawl("https://cs272-f24.github.io/tests/rnj/", indx)

	for {
		//Run the server until manual exit
		time.Sleep(1 * time.Second)
	}
}
