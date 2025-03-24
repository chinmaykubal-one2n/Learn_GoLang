package main

import "fmt"

func main() {
	reverseDigitsOfANumber()
}

func reverseDigitsOfANumber() {
	var number int
	var reversedNumber int

	fmt.Println("Enter a number:")
	fmt.Scanln(&number)

	if number <= 0 {
		fmt.Println("Please enter a valid number")
		return
	}

	for number != 0 {
		reversedNumber = reversedNumber*10 + number%10
		number = number / 10
	}

	fmt.Printf("Reversed number is: %v\n", reversedNumber)
}
