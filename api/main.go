package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var indx Indexes
var indexOnce sync.Once
var crawled bool

func initIndex(indexType string, existingDB bool) {
	indexOnce.Do(func() {
		// Initialize your index here
		// For now, we'll just use a placeholder
		indx = MakeInMemoryIndex()
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering mainHandler")
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "Welcome to the search engine!")
	case "/crawl":
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

	crawled = true
	//Add crawl
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Crawl started"))
}

func main() {
	indexType := "inmem"
	existingDB := false

	initIndex(indexType, existingDB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", mainHandler)
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
