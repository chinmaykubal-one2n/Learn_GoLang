// return only the even and multiples of 5 from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20
// Sample Output: 10, 20

package main

import "fmt"

func evenNumbers(numbers []int) []int {
	var listOfEvenNumbers []int

	for _, number := range numbers {
		if number%2 == 0 {
			listOfEvenNumbers = append(listOfEvenNumbers, number)
		}
	}

	return listOfEvenNumbers
}

func multipleOfFive(listOfEvenNumbers []int) []int {
	var listOfMultipleOfFive []int

	for _, number := range listOfEvenNumbers {
		if number%5 == 0 {
			listOfMultipleOfFive = append(listOfMultipleOfFive, number)
		}
	}

	return listOfMultipleOfFive
}

func main() {
	sample_input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	listOfEvenNumbers := evenNumbers(sample_input)
	listOfMultipleOfFive := multipleOfFive(listOfEvenNumbers)

	fmt.Printf("List of even numbers: %d\n", listOfEvenNumbers)
	fmt.Printf("List of even and multiples of 5 numbers: %d\n", listOfMultipleOfFive)
}
