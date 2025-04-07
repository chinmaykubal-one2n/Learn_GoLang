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
	"fmt"
	"unicode"
)

func IsPalindrome(sentence string) bool {

	var cleaned []rune

	for _, character := range sentence {
		if unicode.IsLetter(character) || unicode.IsDigit(character) {
			character = unicode.ToLower(character)
			cleaned = append(cleaned, character)
		}
	}

	for i := range len(cleaned) / 2 { //new syntax written as per the vscode suggestion
		if cleaned[i] != cleaned[len(cleaned)-1-i] {
			return false
		}
	}
	// fmt.Printf("output %#v", cleanString)
	return true
}

func main() {
	if IsPalindrome("No lemon, no melon") {
		fmt.Println("Given string is Palindrome")
	} else {
		fmt.Println("Given string is not Palindrome")
	}
}
