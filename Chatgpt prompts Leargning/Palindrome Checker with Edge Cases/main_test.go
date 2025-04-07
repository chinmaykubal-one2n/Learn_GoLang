// Implement a function IsPalindrome(s string) bool that returns true if the string is a palindrome. Ignore:
//     Case (A == a)
//     Non-alphanumeric characters (e.g., punctuation, spaces)

// IsPalindrome("A man, a plan, a canal: Panama") == true
// IsPalindrome("race a car") == false
// IsPalindrome("") == true
// IsPalindrome("12321") == true
// IsPalindrome("No lemon, no melon") == true

package main

import (
	"testing"
)

type testCaseType struct {
	input  string
	output bool
}

var testCases = []testCaseType{
	{
		input:  "A man, a plan, a canal: Panama",
		output: true,
	},
	{
		input:  "race a car",
		output: false,
	},
	{
		input:  "",
		output: true,
	},
	{
		input:  "12321",
		output: true,
	},
	{
		input:  "No lemon, no melon",
		output: true,
	},
}

func TestIsPalindrome(t *testing.T) {
	for _, testCase := range testCases {
		output := IsPalindrome(testCase.input)
		if output != testCase.output {
			t.Errorf("For input '%s', expected %#v, but got %#v", testCase.input, testCase.output, output)
		}
	}
}
