package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kljensen/snowball"
)

func Crawl(baseURL string, index Indexes) {
	visitedUrls := make(map[string]bool) //Make a map for all visited urls
	stopWordMap := loadStopWords("stopwords-en.json")
	host, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("URL Parse returned %v", err)
	}
	visitedUrls[baseURL] = true
	hostName := host.Host
	//Read the robots.txt file if it exists
	crawlDelay, dissalowList := loadRobots(hostName)
	chDownload := make(chan string, 100)
	chExtract := make(chan downloadResults, 100)
	var mu sync.Mutex     //Make a mutex for the visited map
	var wg sync.WaitGroup //Waitrgoup to find out when all goroutines have finished
	chDownload <- baseURL //Add the first url
	//Start a goroutine to manage all download results
	go func() {
		for currentUrl := range chDownload {
			allowed := true
			for dissalowedPath := range dissalowList {
				matched, _ := regexp.MatchString(dissalowedPath, currentUrl)
				if matched {
					allowed = false
					break
				}
			}
			if allowed {
				wg.Add(1)
				go Download(currentUrl, chExtract, &wg)
				time.Sleep(time.Duration(crawlDelay) * time.Second)
			}
		}
	}()
	//Start a goroutine to manage all extract results
	go func() {
		for content := range chExtract {
			wg.Add(1)
			go func() {
				defer wg.Done()
				words, hrefs := Extract(content.data)
				currentWords := []string{}
				for _, word := range words {
					if stemmedWord, err := snowball.Stem(word, "english", true); err != nil {
						log.Printf("Snowball error: %v", err)
					} else {
						if _, exists := stopWordMap[stemmedWord]; !exists {
							currentWords = append(currentWords, stemmedWord)
						}
					}
				}
				links := Clean(baseURL, hrefs)
				for _, cleanedURL := range links {
					mu.Lock()
					if !visitedUrls[cleanedURL.String()] && hostName == cleanedURL.Host {
						chDownload <- cleanedURL.String()
						visitedUrls[cleanedURL.String()] = true
					}
					mu.Unlock()
				}
				index.AddToIndex(content.url, currentWords)
			}()
		}
	}()

	//Wait for intial goroutines to spin up and call others
	time.Sleep(2 * time.Second)
	wg.Wait()
	close(chDownload)
	close(chExtract)
	fmt.Printf("All goroutines finished")
}

func loadStopWords(link string) map[string]struct{} {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	fullPath := filepath.Join(cwd, link)
	stopWordsFile, err := os.Open(fullPath)
	if err != nil {
		fmt.Println(err)
	}
	defer stopWordsFile.Close()
	byteValue, _ := io.ReadAll(stopWordsFile)
	var stopWords []string
	if json.Unmarshal(byteValue, &stopWords); err != nil {
		log.Printf("Json unmarshal returned: %v", err)
	}
	stopWordMap := make(map[string]struct{})
	for _, word := range stopWords {
		stopWordMap[word] = struct{}{}
	}
	return stopWordMap
}

func loadRobots(hostName string) (float64, map[string]bool) {
	var crawlDelay float64 = 0.1
	robotsUrl := "http://" + hostName + "/top10/robots.txt"
	dissalowList := make(map[string]bool)
	if res, err := downloadRobots(robotsUrl); err != nil {
		log.Println("No robots file found, continuing standard crawling")
	} else {
		lines := strings.Split(res, "\n")
		currUser := false
		for i := range lines {
			if strings.HasPrefix(lines[i], "User-agent:") {
				if strings.HasPrefix(lines[i], "User-agent: *") {
					currUser = true
				} else {
					currUser = false
				}
			} else if currUser && strings.HasPrefix(lines[i], "Disallow:") {
				filePath := strings.TrimSpace(strings.TrimPrefix(lines[i], "Disallow:"))
				dissalowed := strings.ReplaceAll(filePath, "*", ".*")
				dissalowList[dissalowed] = false
			} else if strings.HasPrefix(lines[i], "Crawl-delay:") {
				delay := strings.TrimSpace(strings.TrimPrefix(lines[i], "Crawl-delay:"))
				i, err := strconv.ParseFloat(delay, 64)
				if err != nil {
					log.Println("robots.txt crawl delay incorrectly formatted")
				} else {
					crawlDelay = float64(i)
				}
			}
		}
	}
	return crawlDelay, dissalowList
}
