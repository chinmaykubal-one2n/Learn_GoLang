package main

import "fmt"

func main() {
	sumOfDigitsOfANumber()
}

func sumOfDigitsOfANumber() {
	var number int
	sumOfDigits := 0

	fmt.Println("Enter a number: ")
	fmt.Scanln(&number)

	if number < 0 || number == 0 {
		fmt.Println("Please enter a valid number")
		return
	}

	for {
		sumOfDigits += number % 10
		number = number / 10

		if number == 0 {
			break
		}
	}
	fmt.Printf("sum of digits is %v\n", sumOfDigits)
}
