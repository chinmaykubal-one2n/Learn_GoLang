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
	finalSlice := finalEvaluate(userCondition, numbers)

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
	var hashMap = make(map[int]bool) // name is so for now to understanding the stuff

	for _, condition := range userCondition {
		if strings.Contains(condition, "even") {
			evenOdds = getEvenOdd("even", inputNumbers)
			tempoSlice = append(tempoSlice, evenOdds...)
		} else if strings.Contains(condition, "odd") {
			evenOdds = getEvenOdd("odd", inputNumbers)
			tempoSlice = append(tempoSlice, evenOdds...)
		}
	}

	for _, condition := range userCondition {
		if strings.Contains(condition, "greater than") {
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			greaterLessNos = getGreaterOrLessThan(parts[0], num, inputNumbers)
			tempoSlice = append(tempoSlice, greaterLessNos...)
		} else if strings.Contains(condition, "less than") {
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			greaterLessNos = getGreaterOrLessThan(parts[0], num, inputNumbers)
			tempoSlice = append(tempoSlice, greaterLessNos...)
		}
	}

	for _, condition := range userCondition {
		if strings.Contains(condition, "multiple of") {
			parts := strings.Split(condition, " ")
			num, _ := strconv.Atoi(parts[2])
			multiplesOf = getMultipleOfN(num, inputNumbers)
			tempoSlice = append(tempoSlice, multiplesOf...)
		}
	}

	for _, condition := range userCondition {
		if strings.Contains(condition, "prime") {
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
