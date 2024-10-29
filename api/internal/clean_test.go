package internal

import (
	"reflect"
	"testing"
)

func TestCleanHref(t *testing.T) {
	var tests = []struct {
		testname string
		hostName string
		hrefs    []string
		want     []string
	}{
		{"cleanTest1", "https://cs272-f24.github.io/",
			[]string{"/", "/help/", "/syllabus/", "https://gobyexample.com/"},
			[]string{"https://cs272-f24.github.io/", "https://cs272-f24.github.io/help/", "https://cs272-f24.github.io/syllabus/", "https://gobyexample.com/"}},
		{"cleanTest2", "https://go.dev/",
			[]string{"/", "/learn/", "/solutions/", "solutions/case-studies/", "https://pkg.go.dev/"},
			[]string{"https://go.dev/", "https://go.dev/learn/", "https://go.dev/solutions/", "https://go.dev/solutions/case-studies/", "https://pkg.go.dev/"}},
		{"cleanTest3", "https://stackoverflow.com/",
			[]string{"/", "/questions/", "/tags/", "/users/", "https://stackoverflow.co/"},
			[]string{"https://stackoverflow.com/", "https://stackoverflow.com/questions/", "https://stackoverflow.com/tags/", "https://stackoverflow.com/users/", "https://stackoverflow.co/"}},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			got := Clean(test.hostName, test.hrefs)
			var gotURLs []string
			for i := range got {
				gotURLs = append(gotURLs, got[i].String())
			}
			if !reflect.DeepEqual(gotURLs, test.want) {
				t.Errorf("\nUrls were: %v\nWanted: %v", gotURLs, test.want)
			}
		})
	}
}
