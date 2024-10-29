package internal

import (
	"reflect"
	"testing"
)

func TestTfIdf(t *testing.T) {
	testCases := []struct {
		name      string
		word      string
		input     InvertedIndex
		documents Frequency
		expected  hits
	}{
		{
			name: "tfIDFTest1",
			word: "romeo",
			input: InvertedIndex{
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
				searchHit{"url6", 10, 0.05555555555555555},
				searchHit{"url1", 4, 0.01},
				searchHit{"url2", 4, 0.01},
			},
		},
		{
			name: "tfIDFTest2",
			word: "juliet",
			input: InvertedIndex{
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
				searchHit{"url7", 3, 0.015},
				searchHit{"url3", 5, 0.008333333333333333},
			},
		},
		{
			name: "tfIDFTest3",
			word: "mercutio",
			input: InvertedIndex{
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
				searchHit{"url5", 7, 0.023333333333333334},
			},
		},
		{
			name: "tfIDFTest4",
			word: "tybalt",
			input: InvertedIndex{
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
				searchHit{"url9", 3, 0.0125},
				searchHit{"url8", 2, 0.010526315789473684},
				searchHit{"url4", 1, 0.004},
			},
		},
		{
			name: "tfIDFTest5",
			word: "benvolio",
			input: InvertedIndex{
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
			idx := MakeInMemoryIndex()
			idx.doclen = test.documents
			idx.wordFreq = test.input
			result := idx.Search(test.word)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Got: %v\nExpected: %v", result, test.expected)
			}
		})
	}
}
