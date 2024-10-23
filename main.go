package main

import (
	"database/sql"
	"flag"
	"log"
	"time"
)

func main() {
	//Flags:
	//indexType: either inmem or runs as a db
	//existingDB run with the name of a db to search on, no indexing will take place
	//link will run the crawler on the link. Will output a file called index.db that can be searched on localhost
	indexType := flag.String("index", "", "Specify the index type")
	existingDB := flag.String("existingDB", "", "Specify an existing database to skip crawling and only run the web server")
	link := flag.String("link", "", "Specify a link to crawl")
	flag.Parse()

	var indx Indexes
	if *indexType == "inmem" {
		indx = MakeInMemoryIndex()
		go webserver(indx)
		if *link != "" {
			crawl(*link, indx) //Crawl the user specified link
		}
	} else if *existingDB != "" {
		// Open existing DB and run only the web server
		db, err := sql.Open("sqlite", *existingDB)
		if err != nil {
			log.Printf("DB open returned %v\n", err)
			return
		}
		defer db.Close()
		indx = MakeDBIndex(db)
		go webserver(indx)
	} else {
		// Open or create a new database and crawl
		db, err := sql.Open("sqlite", "index.db")
		if err != nil {
			log.Printf("DB open returned %v\n", err)
			return
		}
		defer db.Close()
		indx = MakeDBIndex(db)
		go webserver(indx)
		if *link != "" {
			crawl(*link, indx) //Crawl the user specified link
		}
	}

	for {
		//Run the server until manual exit
		time.Sleep(1 * time.Second)
	}
}
