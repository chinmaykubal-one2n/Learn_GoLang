// Write a function LongestWord(sentence string) string that returns the longest word in a given sentence.
// Ignore punctuation.
// Words are separated by spaces.
// If there are multiple longest words with the same length, return the first one.
// An empty sentence should return "".

// LongestWord("The quick brown fox jumps over the lazy dog.") == "quick"
// LongestWord("Hello, world!") == "Hello"
// LongestWord("") == ""
// LongestWord("To be or not to be") == "not"

// Testing Strategy:
// Write table-driven tests with:
// Normal sentences
// Sentences with punctuation
// Sentences with multiple longest words
// Empty input
// Numbers in words (e.g., "Go123")
package main

import "testing"

type testCaseType struct {
	input  string
	output string
}

var testCases = []testCaseType{
	{
		input:  "The quick brown fox jumps over the lazy dog.",
		output: "quick",
	},
	{
		input:  "Hello, world!",
		output: "Hello",
	},
	{
		input:  "",
		output: "",
	},
	{
		input:  "To be or not to be",
		output: "not",
	},
}

func TestLongestWord(t *testing.T) {
	for _, testCase := range testCases {
		output := LongestWord(testCase.input)
		if output != testCase.output {
			t.Errorf("For input %#v Expeted is %#v and got %#v", testCase.input, testCase.output, output)
		}
	}
}
