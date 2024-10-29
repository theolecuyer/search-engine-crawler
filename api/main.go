package main

import (
	"database/sql"
	"flag"
	"log"
	"time"
)

func main() {
	indexType := flag.String("index", "", "Specify the index type")
	existingDB := flag.Bool("existingDB", false, "Use the existing index.db file")
	flag.Parse()

	var indx Indexes
	if *indexType == "inmem" {
		indx = MakeInMemoryIndex()
		go Webserver(indx)
		Crawl("http://www.usfca.edu/", indx)
	} else {
		db, err := sql.Open("sqlite", "index.db")
		if err != nil {
			log.Printf("DB open returned %v\n", err)
		}
		defer db.Close()
		indx = MakeDBIndex(db)
		go Webserver(indx)
		if !*existingDB {
			Crawl("http://www.usfca.edu/", indx)
		}
	}

	for {
		//Run the server until manual exit
		time.Sleep(1 * time.Second)
	}
}
