package main

import (
	"reflect"
	"testing"
)

func TestExtract(t *testing.T) {
	var tests = []struct {
		testname  string
		body      []byte
		wordsWant []string
		hrefsWant []string
	}{
		{"extractTest1", []byte(`
		<body>
		<p>Some text here</p>
		  <a href="http://example.com">Example</a>
		</body>
		`), []string{"Some", "text", "here", "Example"}, []string{"http://example.com"}},
		{"extractTest2", []byte(`
		<html>
		<head>
			<title>CS272 | Welcome</title>
		</head>
		<body>
			<p>Hello World!</p>
			<p>Welcome to <a href="https://cs272-f24.github.io/">CS272</a>!</p>
		</body>
		</html>
		`), []string{"CS272", "Welcome", "Hello", "World", "Welcome", "to", "CS272"}, []string{"https://cs272-f24.github.io/"}},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			wordsGot, hrefsGot := extract(string(test.body))
			if !reflect.DeepEqual(wordsGot, test.wordsWant) || !reflect.DeepEqual(hrefsGot, test.hrefsWant) {
				t.Errorf("\nWords were: %v\nWanted: %v\nHrefs were: %v\nWanted: %v", wordsGot, test.wordsWant, hrefsGot, test.hrefsWant)
			}
		})
	}
}
