package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
)

var indx Indexes
var indexOnce sync.Once
var crawled bool

func initIndex(indexType string, existingDB bool) {
	indexOnce.Do(func() {
		if indexType == "inmem" {
			indx = MakeInMemoryIndex()
		} else {
			db, err := sql.Open("sqlite", "index.db")
			if err != nil {
				log.Fatalf("DB open returned %v\n", err)
			}
			defer db.Close()
			indx = MakeDBIndex(db)
			if !existingDB {
				Crawl("http://www.usfca.edu/", indx) // Start crawling if not using an existing DB
			}
		}
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering mainHandler")
	switch r.URL.Path {
	case "/":
		CrawlHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func CrawlHandler(w http.ResponseWriter, r *http.Request) {
	if crawled {
		http.Error(w, "Crawl already started", http.StatusConflict)
		return
	}

	go func() {
		Crawl("http://www.usfca.edu/", indx)
		crawled = true
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Crawl started"))
}

func main() {
	indexType := "inmem"
	existingDB := false

	initIndex(indexType, existingDB)

	http.HandleFunc("/", mainHandler)
}
