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

import (
	"fmt"
	"strings"
	"unicode"
)

func mostFrequentWord(sentence string) string {
	var mostFrequentWord string = ""
	var cleanString strings.Builder
	var cleanWords []string
	var hashMap = make(map[string]int)
	var maxCount int = 0

	if sentence == "" {
		return mostFrequentWord
	}

	for _, character := range sentence {
		if unicode.IsLetter(character) || unicode.IsSpace(character) {
			character = unicode.ToLower(character)
			cleanString.WriteRune(character)
		}
	}

	// fmt.Println(cleanString.String())
	cleanWords = strings.Fields(cleanString.String())
	// fmt.Println(cleanWords)

	for _, cleanWord := range cleanWords {
		hashMap[cleanWord]++
		if hashMap[cleanWord] > maxCount {
			maxCount = hashMap[cleanWord]
			mostFrequentWord = cleanWord
		}
	}
	// fmt.Println(hashMap)
	return mostFrequentWord
}

func main() {
	fmt.Println(mostFrequentWord("the boss is gg gg gg`.`, `,`, `!`"))
}
