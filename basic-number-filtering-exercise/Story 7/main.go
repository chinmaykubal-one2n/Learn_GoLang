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

// TEST CASES ARE NEEDED .

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	var numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	var userCondition = []string{"greater than 5", "multiple of 3", "even"}
	var finalSlice []int

	if len(userCondition) == 0 {
		fmt.Println("No condition provided.")
		return
	} else if len(numbers) == 0 {
		fmt.Println("No number provided.")
		return
	} else {
		finalSlice = finalEvaluate(userCondition, numbers)
	}

	if len(finalSlice) != 0 {
		fmt.Println(finalSlice)
	} else {
		fmt.Println("No number found.")
	}

}

func finalEvaluate(userCondition []string, inputNumbers []int) []int {
	const (
		EvenFlag     = "even"
		OddFlag      = "odd"
		GreaterFlag  = "greater"
		LessFlag     = "less"
		MultipleFlag = "multiple"
		PrimeFlag    = "prime"
	)
	finalSlice := inputNumbers

	for _, condition := range userCondition {
		switch {
		case strings.Contains(condition, EvenFlag):
			finalSlice = getEvenOdd(EvenFlag, finalSlice)
		case strings.Contains(condition, OddFlag):
			finalSlice = getEvenOdd(OddFlag, finalSlice)
		case strings.Contains(condition, GreaterFlag):
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			finalSlice = getGreaterOrLessThan(parts[0], num, finalSlice)
		case strings.Contains(condition, LessFlag):
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			finalSlice = getGreaterOrLessThan(parts[0], num, finalSlice)
		case strings.Contains(condition, MultipleFlag):
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			finalSlice = getMultipleOfN(num, finalSlice)
		case strings.Contains(condition, PrimeFlag):
			finalSlice = getPrimes(finalSlice)
		}
	}

	return finalSlice
}

func getEvenOdd(flagValue string, numbers []int) []int {
	// by default it is an empty array so if user asks for even
	// numbers and if there are no even numbers
	// then it will return an empty array.
	var evenNumbers []int
	var oddNumbers []int

	// this is a magic value. define it first
	const (
		EvenFlag = "even"
		OddFlag  = "odd"
	)

	for _, number := range numbers {
		if number%2 == 0 {
			evenNumbers = append(evenNumbers, number)
		} else {
			oddNumbers = append(oddNumbers, number)
		}
	}

	// below if statements can be optimzed furthur
	if flagValue == EvenFlag { // magic values should not be there define them first
		return evenNumbers
	}

	if flagValue == OddFlag { // magic values should not be there  define them first
		return oddNumbers
	}

	return []int{} // as we have 2 arrays so better return []int{} + this function needs to written somthing so the line is here
}

func getGreaterOrLessThan(flagValue string, greaterLessThanN int, numbers []int) []int {
	var listOfnumbers []int
	const (
		GreaterFlag = "greater"
		LessFlag    = "less"
	)

	if flagValue == GreaterFlag {
		for _, number := range numbers {
			if number > greaterLessThanN {
				listOfnumbers = append(listOfnumbers, number)
			}
		}
	}

	if flagValue == LessFlag {
		for _, number := range numbers {
			if number < greaterLessThanN {
				listOfnumbers = append(listOfnumbers, number)
			}
		}
	}

	return listOfnumbers
}

func getMultipleOfN(multipleOfN int, numbers []int) []int {
	var multiplesOfN []int

	for _, number := range numbers {
		if number%multipleOfN == 0 {
			multiplesOfN = append(multiplesOfN, number)
		}
	}

	return multiplesOfN
}

func getPrimes(numbers []int) []int {
	var primeNumbers []int

	for _, number := range numbers {
		if checkPrimes(number) {
			primeNumbers = append(primeNumbers, number)
		}
	}

	return primeNumbers // write a common funciton and call it in the respective functions in the first step rather than in the end
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
