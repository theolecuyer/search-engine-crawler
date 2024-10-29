package lib

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		testname string
		expected string
	}{
		{
			"downloadTest1",
			`<html><bodyHello CS 272,/body></html>`,
		},
		{
			"downloadTest2",
			`
			<body>
			<p>Some text here</p>
			  <a href="http://example.com">Example</a>
			</body>
			`,
		}, {
			"downloadTest3",
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
		},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(test.expected))
			}))
			got, _ := downloadRobots(svr.URL)
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.expected)
			}
		})
	}
}
