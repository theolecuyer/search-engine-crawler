package main

import (
	"reflect"
	"testing"
)

func TestTfIdf(t *testing.T) {
	testCases := []struct {
		name      string
		word      string
		input     Index
		documents Frequency
		expected  hits
	}{
		{
			name: "tfIDFTest1",
			word: "romeo",
			input: Index{
				"romeo": Frequency{
					"url2": 4,
					"url1": 4,
					"url6": 10,
				},
			},
			documents: Frequency{
				"url1":  400,
				"url2":  400,
				"url3":  600,
				"url4":  250,
				"url5":  300,
				"url6":  180,
				"url7":  200,
				"url8":  190,
				"url9":  240,
				"url10": 330,
			},
			expected: hits{
				searchHit{"url6", 10, 0.05090504065967528},
				searchHit{"url1", 4, 0.009162907318741552},
				searchHit{"url2", 4, 0.009162907318741552},
			},
		},
		{
			name: "tfIDFTest2",
			word: "juliet",
			input: Index{
				"juliet": Frequency{
					"url3": 5,
					"url7": 3,
				},
			},
			documents: Frequency{
				"url1":  400,
				"url2":  400,
				"url3":  600,
				"url4":  250,
				"url5":  300,
				"url6":  180,
				"url7":  200,
				"url8":  190,
				"url9":  240,
				"url10": 330,
			},
			expected: hits{
				searchHit{"url7", 3, 0.01805959206488904},
				searchHit{"url3", 5, 0.010033106702716134},
			},
		},
		{
			name: "tfIDFTest3",
			word: "mercutio",
			input: Index{
				"mercutio": Frequency{
					"url5": 7,
				},
			},
			documents: Frequency{
				"url1":  400,
				"url2":  400,
				"url3":  600,
				"url4":  250,
				"url5":  300,
				"url6":  180,
				"url7":  200,
				"url8":  190,
				"url9":  240,
				"url10": 330,
			},
			expected: hits{
				searchHit{"url5", 7, 0.03755355129012901},
			},
		},
		{
			name: "tfIDFTest4",
			word: "tybalt",
			input: Index{
				"tybalt": Frequency{
					"url4": 1,
					"url8": 2,
					"url9": 3,
				},
			},
			documents: Frequency{
				"url1":  400,
				"url2":  400,
				"url3":  600,
				"url4":  250,
				"url5":  300,
				"url6":  180,
				"url7":  200,
				"url8":  190,
				"url9":  240,
				"url10": 330,
			},
			expected: hits{
				searchHit{"url9", 3, 0.01145363414842694},
				searchHit{"url8", 2, 0.009645165598675317},
				searchHit{"url4", 1, 0.0036651629274966203},
			},
		},
		{
			name: "tfIDFTest5",
			word: "benvolio",
			input: Index{
				"benvolio": Frequency{},
			},
			documents: Frequency{
				"url1":  400,
				"url2":  400,
				"url3":  600,
				"url4":  250,
				"url5":  300,
				"url6":  180,
				"url7":  200,
				"url8":  190,
				"url9":  240,
				"url10": 330,
			},
			expected: hits{},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := search(test.word, test.input, test.documents)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Got: %v\nExpected: %v", result, test.expected)
			}
		})
	}
}
