package main

import "fmt"

func main() {
	message := evenOrOdd()
	fmt.Println(message)
}

func evenOrOdd() string {
	var message string
	var number int

	fmt.Println("Enter a number:")
	fmt.Scan(&number)

	if number%2 == 0 {
		message = "The number is even"
	} else {
		message = "The number is odd"
	}

	return message
}
