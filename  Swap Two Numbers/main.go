package main

import "fmt"

func main() {
	swapTwoNumbers()
}

func swapTwoNumbers() {
	var number_one int
	var number_two int
	var temp_var int

	fmt.Println("Enter the first number:")
	fmt.Scanln(&number_one)

	fmt.Println("Enter the second number:")
	fmt.Scanln(&number_two)

	fmt.Printf("Before swapping, numbers are: %v %v\n", number_one, number_two)

	// swapping logic
	temp_var = number_one
	number_one = number_two
	number_two = temp_var

	fmt.Printf("After swapping, numbers are: %v %v\n", number_one, number_two)
}
