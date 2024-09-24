package main

import (
	"log"
	"net/url"
)

func clean(host string, hrefs []string) []*url.URL {
	var cleaned []*url.URL
	u, err := url.Parse(host)
	if err != nil {
		log.Fatalf("url.Parse returned %v\n", err)
	}
	for i := range hrefs {
		hr, err := url.Parse(hrefs[i])
		if err != nil {
			log.Fatalf("url.Parse returned %v\n", err)
		}
		//If the hrefs is absolute, or its own valid URL, resolve reference will handle the case
		//Used net/url documentation to find the resolve reference method
		newUrl := u.ResolveReference(hr)
		cleaned = append(cleaned, newUrl)
	}
	return cleaned
}
