// Given a list of integers, and a set of conditions (odd, even, greater than 5, multiple of 3, prime,
// and many more such custom conditions that may be dynamically defined by user),
// write a program to return only the integers from the given list that match ALL the conditions.

// Sample Input:
// A list containing 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20
// Conditions specified using a set of functions: odd, greater than 5, multiple of 3
// Sample Output: 9, 15

// Sample Input:
// A list containing 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20
// Conditions specified using a set of functions:

// Sample Output: 6, 12

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	var numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	var userCondition = []string{"greater than 5", "multiple of 3", "odd"}

	finalSlice := finalEvaluate(userCondition, numbers)

	fmt.Println(finalSlice)

}

func finalEvaluate(userCondition []string, inputNumbers []int) []int {
	var finalSlice []int = inputNumbers

	for _, condition := range userCondition {
		if strings.Contains(condition, "even") {
			finalSlice = getEvenOdd("even", finalSlice)
		} else if strings.Contains(condition, "odd") {
			finalSlice = getEvenOdd("odd", finalSlice)
		}
	}

	for _, condition := range userCondition {
		if strings.Contains(condition, "greater than") {
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			finalSlice = getGreaterOrLessThan(parts[0], num, finalSlice)
		} else if strings.Contains(condition, "less than") {
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			finalSlice = getGreaterOrLessThan(parts[0], num, finalSlice)
		}
	}

	for _, condition := range userCondition {
		if strings.Contains(condition, "multiple of") {
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			finalSlice = getMultipleOfN(num, finalSlice)
		}
	}

	for _, condition := range userCondition {
		if strings.Contains(condition, "prime") {
			finalSlice = getPrimes(finalSlice)
		}
	}

	if len(finalSlice) != 0 {
		return finalSlice
	}

	fmt.Println("No number found.")
	return nil
}

func getEvenOdd(flagValue string, numbers []int) []int {
	var evenNumbers []int
	var oddNumbers []int

	for _, number := range numbers {
		if number%2 == 0 {
			evenNumbers = append(evenNumbers, number)
		} else {
			oddNumbers = append(oddNumbers, number)
		}
	}

	if flagValue == "even" {
		return evenNumbers
	}
	if flagValue == "odd" {
		return oddNumbers
	}

	return nil
}

func getGreaterOrLessThan(flagValue string, greaterLessThanN int, numbers []int) []int {

	var listOfnumbers []int

	if flagValue == "greater" {
		for _, number := range numbers {
			if number > greaterLessThanN {
				listOfnumbers = append(listOfnumbers, number)
			}
		}
	}

	if flagValue == "less" {
		for _, number := range numbers {
			if number < greaterLessThanN {
				listOfnumbers = append(listOfnumbers, number)
			}
		}
	}

	if len(listOfnumbers) != 0 {
		return listOfnumbers
	}

	return nil
}

func getMultipleOfN(multipleOfN int, numbers []int) []int {
	var multiplesOfN []int

	for _, number := range numbers {
		if number%multipleOfN == 0 {
			multiplesOfN = append(multiplesOfN, number)
		}
	}

	if len(multiplesOfN) != 0 {
		return multiplesOfN
	}

	return nil
}

func getPrimes(numbers []int) []int {
	var primeNumbers []int

	for _, number := range numbers {
		if checkPrimes(number) {
			primeNumbers = append(primeNumbers, number)
		}
	}

	if len(primeNumbers) != 0 {
		return primeNumbers
	}

	return nil
}

func checkPrimes(number int) bool {

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
