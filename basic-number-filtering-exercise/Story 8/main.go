// story 8
// write a program to return only the integers from the given list that match ANY of the conditions.

// Sample Input:
// A list containing 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20
// Conditions specified using a set of functions: prime, greater than 15, multiple of 5
// Sample Output: 2, 3, 5, 7, 10, 11, 13, 15, 16, 17, 18, 19, 20

// Sample Input:
// A list containing 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20
// Conditions specified using a set of functions: less than 6, multiple of 3
// Sample Output: 1, 2, 3, 4, 5, 6, 9, 12, 15, 18

package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	var userCondition = []string{"less than 6", "odd", "multiple of 3"}
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
	var evenOdds []int = inputNumbers
	var greaterLessNos []int = inputNumbers
	var multiplesOf []int = inputNumbers
	var primes []int = inputNumbers
	var tempoSlice []int
	var finalSlice []int
	const (
		EvenFlag     = "even"
		OddFlag      = "odd"
		GreaterFlag  = "greater"
		LessFlag     = "less"
		MultipleFlag = "multiple"
		PrimeFlag    = "prime"
	)
	var hashMap = make(map[int]bool) // name is so for now to understanding the stuff

	for _, condition := range userCondition {
		switch {
		case strings.Contains(condition, EvenFlag):
			evenOdds = getEvenOdd(EvenFlag, inputNumbers)
			tempoSlice = append(tempoSlice, evenOdds...)
		case strings.Contains(condition, OddFlag):
			evenOdds = getEvenOdd(OddFlag, inputNumbers)
			tempoSlice = append(tempoSlice, evenOdds...)
		case strings.Contains(condition, GreaterFlag):
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			greaterLessNos = getGreaterOrLessThan(parts[0], num, inputNumbers)
			tempoSlice = append(tempoSlice, greaterLessNos...)
		case strings.Contains(condition, LessFlag):
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			greaterLessNos = getGreaterOrLessThan(parts[0], num, inputNumbers)
			tempoSlice = append(tempoSlice, greaterLessNos...)
		case strings.Contains(condition, MultipleFlag):
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			multiplesOf = getMultipleOfN(num, inputNumbers)
			tempoSlice = append(tempoSlice, multiplesOf...)
		case strings.Contains(condition, PrimeFlag):
			primes = getPrimes(inputNumbers)
			tempoSlice = append(tempoSlice, primes...)
		}
	}

	sort.Ints(tempoSlice)

	for _, number := range tempoSlice {
		if !hashMap[number] {
			hashMap[number] = true
			finalSlice = append(finalSlice, number)
		}
	}

	return finalSlice
}

func getEvenOdd(flagValue string, numbers []int) []int {
	var evenNumbers []int
	var oddNumbers []int
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

	if flagValue == EvenFlag {
		return evenNumbers
	}
	if flagValue == OddFlag {
		return oddNumbers
	}

	return []int{}
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

	return primeNumbers
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
