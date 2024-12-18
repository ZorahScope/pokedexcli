package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " weekend warrior          ",
			expected: []string{"weekend", "warrior"},
		},
		{
			input:    "Grand Burrito",
			expected: []string{"grand", "burrito"},
		},
		{
			input:    "Cows And Purple muffins",
			expected: []string{"cows", "and", "purple", "muffins"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		actualSliceLength := len(actual)
		expectedSliceLength := len(c.expected)
		if actualSliceLength != expectedSliceLength {
			t.Errorf("slice length %v from cleanInput(%#v) != expected %v", actualSliceLength, actual, expectedSliceLength)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%v) returned words not matching whats expected", c.input)
			}
		}
	}
}
