// return only the odd numbers from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
// Sample Output: 1, 3, 5, 7, 9
package main

import "fmt"

func main() {
	var sampleInput = []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	fmt.Printf("The list of Odd numbers: %d\n", oddNumbers(sampleInput))
}

func oddNumbers(sampleInput []int) []int {
	var listOfOddNumbers []int

	for _, number := range sampleInput {
		if number%2 != 0 {
			listOfOddNumbers = append(listOfOddNumbers, number)
		}
	}

	return listOfOddNumbers
}
