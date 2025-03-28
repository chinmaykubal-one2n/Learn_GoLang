// return only the odd and multiples of 3 greater than 10 from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20
// Sample Output: 15

package main

import "fmt"

func getOddNumbers(numbers []int) []int {
	var oddNumbers []int

	for _, number := range numbers {
		if number%2 != 0 {
			oddNumbers = append(oddNumbers, number)
		}
	}

	return oddNumbers
}

func getMultiplesOfThree(numbers []int) []int {
	var multiplesOfThree []int

	for _, number := range numbers {
		if number%3 == 0 {
			multiplesOfThree = append(multiplesOfThree, number)
		}
	}

	return multiplesOfThree
}

func getGreaterThanTen(numbers []int) []int {
	var numbersGreaterThanTen []int

	for _, number := range numbers {
		if number > 10 {
			numbersGreaterThanTen = append(numbersGreaterThanTen, number)
		}
	}

	return numbersGreaterThanTen
}

func main() {
	var sampleInput = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	oddNumbers := getOddNumbers(sampleInput)
	multiplesOfThree := getMultiplesOfThree(oddNumbers)
	numbersGreaterThanTen := getGreaterThanTen(multiplesOfThree)

	fmt.Printf("List of odd numbers is: %d\n", oddNumbers)
	fmt.Printf("List of odd numbers and multiples of 3 is: %d\n", multiplesOfThree)
	fmt.Printf("List of odd numbers and multiples of 3 and greater than 10 is: %d\n", numbersGreaterThanTen)

}
