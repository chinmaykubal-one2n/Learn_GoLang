// Write a function LongestWord(sentence string) string that returns the longest word in a given sentence.
// Ignore punctuation.
// Words are separated by spaces.
// If there are multiple longest words with the same length, return the first one.
// An empty sentence should return "".

// LongestWord("The quick brown fox jumps over the lazy dog.") == "quick"
// LongestWord("Hello, world!") == "Hello"
// LongestWord("") == ""
// LongestWord("To be or not to be") == "not"

package main

import (
	"fmt"
	"strings"
	"unicode"
)

func LongestWord(sentence string) string {
	var longestWord string = ""
	var cleanString strings.Builder
	var cleanWords []string

	if sentence == "" {
		return longestWord
	}

	for _, character := range sentence {
		if unicode.IsLetter(character) || unicode.IsDigit(character) || unicode.IsSpace(character) {
			cleanString.WriteRune(character)
		}
	}

	// fmt.Printf("cleanString %#v\n", cleanString.String())
	cleanWords = strings.Fields(cleanString.String())
	// fmt.Printf("words %#v\n", words)

	for _, word := range cleanWords {
		if len(word) > len(longestWord) {
			longestWord = word
		}
	}

	return longestWord
}

func main() {
	fmt.Println(LongestWord(""))
}
