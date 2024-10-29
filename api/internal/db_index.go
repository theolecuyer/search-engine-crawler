package internal

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/kljensen/snowball"
	_ "github.com/lib/pq"
)

type DatabaseIndex struct {
	db                  *sql.DB
	sessionID           string
	insertURLStmt       *sql.Stmt
	insertWordStmt      *sql.Stmt
	insertFreqStmt      *sql.Stmt
	getWordIDStmt       *sql.Stmt
	getURLStmt          *sql.Stmt
	getURLWordCountStmt *sql.Stmt
	getWordFreqStmt     *sql.Stmt
}

func MakeDBIndex(db *sql.DB, sessionID string) *DatabaseIndex {
	insertURLStmt := prepare(db, `INSERT INTO urls (url, word_count, user_id) VALUES ($1, $2, (SELECT user_id FROM users WHERE session_id = $3)) ON CONFLICT (url) DO NOTHING`)
	insertWordStmt := prepare(db, `INSERT INTO words (word, user_id) VALUES ($1, (SELECT user_id FROM users WHERE session_id = $2)) ON CONFLICT (word) DO NOTHING`)
	insertFreqStmt := prepare(db, `INSERT INTO mapping (word_id, url_id, frequency) VALUES ($1, $2, 1) ON CONFLICT (word_id, url_id) DO UPDATE SET frequency = mapping.frequency + 1;`)

	//Queries
	getWordIDStmt := prepare(db, `SELECT id FROM words WHERE word = $1 AND user_id = (SELECT user_id FROM users WHERE session_id = $2)`)
	getURLStmt := prepare(db, `SELECT url FROM urls WHERE id = $1 AND user_id = (SELECT user_id FROM users WHERE session_id = $2)`)
	getURLWordCountStmt := prepare(db, `SELECT word_count FROM urls WHERE id = $1 AND user_id = (SELECT user_id FROM users WHERE session_id = $2)`)
	getWordFreqStmt := prepare(db, `SELECT url_id, frequency FROM mapping WHERE word_id = $1`)
	return &DatabaseIndex{
		db:                  db,
		sessionID:           sessionID,
		insertURLStmt:       insertURLStmt,
		insertWordStmt:      insertWordStmt,
		insertFreqStmt:      insertFreqStmt,
		getWordIDStmt:       getWordIDStmt,
		getURLStmt:          getURLStmt,
		getURLWordCountStmt: getURLWordCountStmt,
		getWordFreqStmt:     getWordFreqStmt,
	}
}

func (d *DatabaseIndex) AddToIndex(url string, currWords []string) {
	//Use transactions to batch apply queries to the db
	tx, err := d.db.Begin()
	if err != nil {
		log.Printf("Tx returned: %v\n", err)
		return
	}
	fmt.Printf("Indexing DB at: %s\n", url)

	// Insert the URL with session ID
	res, err := tx.Stmt(d.insertURLStmt).Exec(url, len(currWords), d.sessionID)
	if err != nil {
		log.Printf("URL insert returned %v\n", err)
		tx.Rollback()
		return
	}

	urlID, err := res.LastInsertId()
	if err != nil {
		log.Printf("URL last insert returned %v\n", err)
		tx.Rollback()
		return
	}

	for _, word := range currWords {
		var wordID int64
		// Insert the word with session ID
		res, err := tx.Stmt(d.insertWordStmt).Exec(word, d.sessionID)
		if err != nil {
			log.Printf("Word insert returned %v\n", err)
			tx.Rollback()
			return
		}

		num, err := res.RowsAffected()
		if err != nil {
			log.Printf("Rows affected returned %v\n", err)
			tx.Rollback()
			return
		}

		if num != 0 {
			wordID, err = res.LastInsertId()
			if err != nil {
				log.Printf("Word insert returned %v\n", err)
				tx.Rollback()
				return
			}
		} else {
			err := tx.Stmt(d.getWordIDStmt).QueryRow(word).Scan(&wordID)
			if err != nil {
				log.Printf("Get word ID returned %v\n", err)
				tx.Rollback()
				return
			}
		}

		_, err = tx.Stmt(d.insertFreqStmt).Exec(wordID, urlID)
		if err != nil {
			log.Printf("Insert frequency returned %v\n", err)
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Tx commit returned %v\n", err)
	}
}

func (d *DatabaseIndex) Search(query string) hits {
	results := hits{}
	if stemmedWordQuery, err := snowball.Stem(query, "english", true); err == nil {
		wordId := d.getWordID(stemmedWordQuery)
		rows, err := d.getWordFreqStmt.Query(wordId)
		if err != nil {
			log.Panicf("Lookup word freq returned %v\n", err)
		}
		defer rows.Close()
		resultUrl := make(map[int]int) // URL id: frequency
		for rows.Next() {
			var wordFreq int
			var urlID int
			err := rows.Scan(&urlID, &wordFreq)
			if err != nil {
				log.Printf("Failed to scan row %v\n", err)
			}
			resultUrl[urlID] = wordFreq
		}
		row := d.db.QueryRow("SELECT COUNT(*) FROM urls WHERE user_id = (SELECT user_id FROM users WHERE session_id = $1)", d.sessionID)
		var totalDocCount int
		if err := row.Scan(&totalDocCount); err != nil {
			log.Printf("Error counting rows: %v", err)
		}
		for url, frequency := range resultUrl {
			currURL := d.getURL(url)
			docLen := d.getURLWordCount(url)
			tfIDFScore := TfIDF(frequency, docLen, totalDocCount, len(resultUrl))
			results = append(results, searchHit{currURL, frequency, tfIDFScore})
		}
	}
	sort.Sort(results)
	return results
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
	err := d.getWordIDStmt.QueryRow(word, d.sessionID).Scan(&wordID)
	if err != nil {
		log.Printf("Lookup for %s returned %v\n", word, err)
	}
	return wordID
}

func (d *DatabaseIndex) getURL(urlID int) string {
	var url string
	err := d.getURLStmt.QueryRow(urlID, d.sessionID).Scan(&url)
	if err != nil {
		log.Printf("Lookup for URL ID %d returned %v\n", urlID, err)
	}
	return url
}

func (d *DatabaseIndex) getURLWordCount(urlID int) int {
	var wordCount int
	err := d.getURLWordCountStmt.QueryRow(urlID, d.sessionID).Scan(&wordCount)
	if err != nil {
		log.Printf("Lookup for URL ID %d returned %v\n", urlID, err)
	}
	return wordCount
}
