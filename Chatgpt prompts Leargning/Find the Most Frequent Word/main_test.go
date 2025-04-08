//  **Challenge: Find the Most Frequent Word**

// Write a function that:
// - Takes a sentence as input.
// - Returns the **most frequent word** (case-insensitive).
// - Ignores punctuation.

// Input: "The fox jumped over the dog. The dog barked."
// Output: "the"

// - Words should be compared in **lowercase**.
// - Ignore punctuation like `.`, `,`, `!`, etc.
// - In case of a tie, return the first word that reached the highest count.

package main

import "testing"

type testCaseType struct {
	input  string
	output string
}

var testCases = []testCaseType{
	{
		input:  "The fox jumped over the dog. The dog barked.",
		output: "the",
	},
	{
		input:  "The fox jumped over the dog, dog!! The dog barked.", // In case of a tie, return the first word that reached the highest count.
		output: "the",
	},
}

func TestMostFrequentWord(t *testing.T) {
	for _, testCase := range testCases {
		output := mostFrequentWord(testCase.input)
		if output != testCase.output {
			t.Errorf("For input %#v expected is %#v and got %#v", testCase.input, testCase.output, output)
		}
	}
}
