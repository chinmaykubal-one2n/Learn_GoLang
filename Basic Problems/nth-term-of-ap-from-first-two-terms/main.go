package main

import "fmt"

func main() {
	nthTermOfAPFromFirstTwoTerms()
}

func nthTermOfAPFromFirstTwoTerms() {
	var first_number int
	var second_number int
	var nth_term int
	var nth_term_value int
	var difference int
	var next_number int
	arithmatic_progression := []int{}

	fmt.Println("Enter the first number of the AP: ")
	fmt.Scanln(&first_number)
	next_number = first_number

	fmt.Println("Enter the second number of the AP: ")
	fmt.Scanln(&second_number)

	fmt.Println("Enter the nth term of the AP: ")
	fmt.Scanln(&nth_term)

	difference = second_number - first_number

	arithmatic_progression = append(arithmatic_progression, next_number)
	// fmt.Println(next_number)

	for i := 0; i < nth_term-1; i++ {
		next_number = next_number + difference
		arithmatic_progression = append(arithmatic_progression, next_number)
		// fmt.Println(next_number)
	}

	nth_term_value = next_number
	fmt.Println("The AP is:", arithmatic_progression)
	fmt.Println("The nth term of the AP is: ", nth_term_value)

}
