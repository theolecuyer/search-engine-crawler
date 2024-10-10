package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/kljensen/snowball"
	_ "modernc.org/sqlite"
)

type DatabaseIndex struct {
	db                  *sql.DB
	insertURLStmt       *sql.Stmt
	insertWordStmt      *sql.Stmt
	insertFreqStmt      *sql.Stmt
	getWordIDStmt       *sql.Stmt
	getURLStmt          *sql.Stmt
	getURLWordCountStmt *sql.Stmt
	getWordFreqStmt     *sql.Stmt
}

func MakeDBIndex(db *sql.DB) *DatabaseIndex {
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Printf("Foreign key error: %v\n", err)
	}
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		log.Fatalf("Failed to set WAL mode: %v", err)
	}
	tables := []string{
		`CREATE TABLE IF NOT EXISTS urls(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT UNIQUE NOT NULL,
		word_count INTEGER NOT NULL
	);`,
		`CREATE TABLE IF NOT EXISTS words(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        word TEXT UNIQUE NOT NULL
	);`,
		`CREATE TABLE IF NOT EXISTS mapping(
		word_id INTEGER,
        url_id INTEGER,
		frequency INTEGER NOT NULL,
		FOREIGN KEY(word_id) REFERENCES words(id),
		FOREIGN KEY(url_id) REFERENCES urls(id),
		UNIQUE (word_id, url_id)
	);`,
	}
	createTables(db, tables)
	insertURLStmt := prepare(db, `INSERT OR IGNORE INTO urls (url, word_count) VALUES (?, ?)`)
	insertWordStmt := prepare(db, `INSERT OR IGNORE INTO words (word) VALUES (?)`)
	insertFreqStmt := prepare(db, `INSERT INTO mapping (word_id, url_id, frequency) VALUES (?, ?, 1) ON CONFLICT(word_id, url_id) DO UPDATE SET frequency = mapping.frequency + 1;`)
	getWordIDStmt := prepare(db, `SELECT id FROM words WHERE word = ?`)
	getURLStmt := prepare(db, `SELECT url FROM urls WHERE id = ?`)
	getURLWordCountStmt := prepare(db, `SELECT word_count FROM urls WHERE id = ?`)
	getWordFreqStmt := prepare(db, `SELECT url_id, frequency FROM mapping WHERE word_id = ?`)
	return &DatabaseIndex{
		db:                  db,
		insertURLStmt:       insertURLStmt,
		insertWordStmt:      insertWordStmt,
		insertFreqStmt:      insertFreqStmt,
		getWordIDStmt:       getWordIDStmt,
		getURLStmt:          getURLStmt,
		getURLWordCountStmt: getURLWordCountStmt,
		getWordFreqStmt:     getWordFreqStmt,
	}
}

func (d *DatabaseIndex) AddToIndex(allWords map[string][]string) {
	//Use transactions to batch apply queries to the db
	tx, err := d.db.Begin()
	if err != nil {
		log.Printf("Tx returned: %v\n", err)
	}
	for url, words := range allWords {
		fmt.Printf("Indexing DB at: %s\n", url)
		res, err := tx.Stmt(d.insertURLStmt).Exec(url, len(words))
		if err != nil {
			log.Printf("URL insert returned %v\n", err)
			tx.Rollback()
		}
		urlID, err := res.LastInsertId()
		if err != nil {
			log.Printf("URL last insert returned %v\n", err)
		}
		for _, word := range words {
			var wordID int64
			res, err := tx.Stmt(d.insertWordStmt).Exec(word)
			if err != nil {
				log.Printf("word insert returned %v\n", err)
				tx.Rollback()
			}
			num, err := res.RowsAffected()
			if err != nil {
				log.Printf("Rows affected returned %v\n", err)
				tx.Rollback()
			}
			if num != 0 {
				wordID, err = res.LastInsertId()
				if err != nil {
					log.Printf("URL insert returned %v\n", err)
					tx.Rollback()
				}
			} else {
				err := tx.Stmt(d.getWordIDStmt).QueryRow(word).Scan(&wordID)
				if err != nil {
					log.Printf("Get word ID returned %v\n", err)
					tx.Rollback()
				}
			}
			_, err = tx.Stmt(d.insertFreqStmt).Exec(wordID, urlID)
			if err != nil {
				log.Printf("Insert freq returned %v\n", err)
				tx.Rollback()
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("tx commit returned %v\n", err)
	}
}

func (d *DatabaseIndex) Search(query string) hits {
	results := hits{}
	//resultingURLFreq = make(map[string]int)
	if stemmedWordQuery, err := snowball.Stem(query, "english", true); err == nil {
		wordId := d.getWordID(stemmedWordQuery)
		rows, err := d.getWordFreqStmt.Query(wordId)
		if err != nil {
			log.Panicf("Lookup word freq returned %v\n", err)
		}
		defer rows.Close()
		resultUrl := make(map[int]int) //Url id: frequency
		for rows.Next() {
			var wordFreq int
			var url_id int
			err := rows.Scan(&url_id, &wordFreq)
			if err != nil {
				log.Printf("Failed to scan row %v\n", err)
			}
			resultUrl[url_id] = wordFreq
		}
		row := d.db.QueryRow("SELECT COUNT(*) FROM urls")
		var totalDocCount int
		if err := row.Scan(&totalDocCount); err != nil {
			log.Printf("Error counting rows: %v", err)
		}
		for url, frequency := range resultUrl {
			currURL := d.getURL(url)
			docLen := d.getURLWordCount(url)
			tfIDFScore := tfIDF(frequency, docLen, totalDocCount, len(resultUrl))
			results = append(results, searchHit{currURL, frequency, tfIDFScore})
		}
	}
	sort.Sort(results)
	return results
}

func createTables(db *sql.DB, tables []string) {
	for _, query := range tables {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Table error: %v", err)
		}
	}
}

func prepare(db *sql.DB, statement string) *sql.Stmt {
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Printf("Prepare returned %v\n", err)
	}
	return stmt
}

func (d *DatabaseIndex) getWordID(word string) int {
	var wordID int
	err := d.getWordIDStmt.QueryRow(word).Scan(&wordID)
	if err != nil {
		log.Printf("Lookup for %s returned %v\n", word, err)
	}
	return wordID
}

func (d *DatabaseIndex) getURL(urlID int) string {
	var url string
	err := d.getURLStmt.QueryRow(urlID).Scan(&url)
	if err != nil {
		log.Printf("Lookup for %s returned %v\n", url, err)
	}
	return url
}

func (d *DatabaseIndex) getURLWordCount(urlID int) int {
	var wordCount int
	err := d.getURLWordCountStmt.QueryRow(urlID).Scan(&wordCount)
	if err != nil {
		log.Printf("Lookup for %d returned %v\n", urlID, err)
	}
	return wordCount
}
