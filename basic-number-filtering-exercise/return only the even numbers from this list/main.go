// write a program to return only the even numbers from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
// Sample Output: 2, 4, 6, 8, 10

package main

import "fmt"

func main() {
	var sampleInput = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("List of even numbers %d\n", evenNumbers(sampleInput))
}

func evenNumbers(listOfNumbers []int) []int {
	var evenNumbers []int

	for _, number := range listOfNumbers {
		if number%2 == 0 {
			evenNumbers = append(evenNumbers, number)
		}
	}

	return evenNumbers
}
