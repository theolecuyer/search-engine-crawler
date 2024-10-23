package main

import (
	"database/sql"
	"flag"
	"log"
	"time"
)

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
	crawl("https://gihthub.com", indx)

	for {
		//Run the server until manual exit
		time.Sleep(1 * time.Second)
	}
}
