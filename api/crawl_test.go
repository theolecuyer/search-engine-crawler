package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sort"
	"testing"
)

func TestCrawl(t *testing.T) {
	tests := []struct {
		testname string
		body     string
		expected []string
	}{
		{
			"crawlTest1",
			`<html>
			<body>
			  <ul>
				<li>
				  <a href="/simple.html">simple.html</a>
				</li>
				<li>
				  <a href="/href.html">href.html</a>
				</li>
				<li>
				  <a href="/style.html">style.html</a>
			  </ul>
			</body>
			</html>`, []string{"/", "/simple.html", "/href.html", "/style.html"},
		},
		{
			"crawlTest2",
			`<html>
			<body>
			Hello CS 272, there are no links here.
			</body>
			</html>`, []string{"/"},
		},
		{
			"crawlTest3",
			`<html>
			<body>
			For a simple example, see <a href="/simple.html">simple.html</a>
			</body>
			</html>`, []string{"/", "/simple.html"},
		},
		{
			"crawlTest4",
			`<html>
			<head>
			  <title>Style</title>
			  <style>
				a.blue {
				  color: blue;
				}
				a.red {
				  color: red;
				}
			  </style>
			<body>
			  <p>
				Here is a blue link to <a class="blue" href="/href.html">href.html</a>
			  </p>
			  <p>
				And a red link to <a class="red" href="/simple.html">simple.html</a>
			  </p>
			</body>
			</html>`, []string{"/", "/href.html", "/simple.html"},
		},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/":
					w.Write([]byte(test.body))
				case "/simple.html":
					w.Write([]byte(tests[1].body))
				case "/href.html":
					w.Write([]byte(tests[2].body))
				case "/style.html":
					w.Write([]byte(tests[3].body))
				default:
					http.NotFound(w, r)
				}
			}))
			defer svr.Close()
			idx := MakeInMemoryIndex()
			Crawl(svr.URL, idx)
			var gotURLS []string
			for u := range idx.doclen {
				parsedUrl, err := url.Parse(u)
				if err != nil {
					t.Fatalf("Failed to parse URL: %v", err)
				}
				if parsedUrl.Path == "" {
					gotURLS = append(gotURLS, "/")
				} else {
					gotURLS = append(gotURLS, parsedUrl.Path)
				}
			}

			// Sort both slices before comparison
			sort.Strings(gotURLS)
			sort.Strings(test.expected)
			if !reflect.DeepEqual(gotURLS, test.expected) {
				t.Errorf("\nGot: %v\nExpected: %v\n", gotURLS, test.expected)
			}
		})
	}
}
