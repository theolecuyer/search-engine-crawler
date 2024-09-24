package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestStop(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]struct{}
	}{
		{
			"stopWordTest1",
			`<html><bodyHello CS 272,/body></html>`,
			map[string]struct{}{},
		},
		{
			"stopWordTest2",
			`
			<body>
			<p>Some text here</p>
			  <a href="http://example.com">Example</a>
			</body>
			`,
			map[string]struct{}{
				"exampl": {},
			},
		}, {
			"stopWordTest3",
			`
		<html>
		<head>
			<title>CS272 | Welcome</title>
		</head>
		<body>
			<p>Hello World!</p>
			<p>Welcome to <a href="https://cs272-f24.github.io/">CS272</a>!</p>
		</body>
		</html>
		`,
			map[string]struct{}{
				"cs272":  {},
				"welcom": {},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(test.input))
			}))
			dataWordFreq := make(Index)
			dataDocLen := make(Frequency)
			results := make(map[string]struct{})
			index(svr.URL, dataWordFreq, dataDocLen)
			for word := range dataWordFreq {
				results[word] = struct{}{}
			}
			if !reflect.DeepEqual(results, test.expected) {
				t.Errorf("Got: %v\nExpected: %v", results, test.expected)
			}
		})
	}
}
