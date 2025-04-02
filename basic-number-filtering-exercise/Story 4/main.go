// return only the odd prime numbers from this list
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
// Sample Output: 3, 5, 7

package main

import (
	"fmt"
	"math"
)

func primeNumbers(number int) bool {
	if number < 2 {
		return false
	}

	for i := 2; i <= int(math.Sqrt(float64(number))); i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}

func oddNumbers(primeNumbers []int) []int {
	var listOfOddNumbers []int

	for _, number := range primeNumbers {

		if number%2 != 0 {
			listOfOddNumbers = append(listOfOddNumbers, number)
		}
	}

	return listOfOddNumbers
}

func main() {
	sample_input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var listOfPrimeNumbers []int

	for _, number := range sample_input {
		if primeNumbers(number) {
			listOfPrimeNumbers = append(listOfPrimeNumbers, number)
		}
	}

	fmt.Printf("Prime numbers are:- %d\n", listOfPrimeNumbers)
	fmt.Printf("Odd Prime numbers are:- %d\n", oddNumbers(listOfPrimeNumbers))
}
