package main

import (
	"fmt"
)

func main() {
	findClosestNumber()
}

func findClosestNumber() {
	var number_to_divide int
	var divisible_by int

	var maximum_value int
	var new_number_to_divide_right int
	var new_number_to_divide_left int
	var right_shift_positions int
	var left_shift_positions int

	fmt.Println("Enter a number:")
	fmt.Scanln(&number_to_divide)

	fmt.Println("Enter a number to divide by:")
	fmt.Scanln(&divisible_by)

	number_to_divide = getAbsoluteValue(number_to_divide)
	divisible_by = getAbsoluteValue(divisible_by)

	if isCompletelyDivisible(number_to_divide, divisible_by) {
		fmt.Printf("The number %v is divisible by %v\n", number_to_divide, divisible_by)
		return
	}

	new_number_to_divide_right = number_to_divide + 1
	for {
		if new_number_to_divide_right%divisible_by == 0 {
			break
		}
		new_number_to_divide_right++
	}

	new_number_to_divide_left = number_to_divide - 1
	for {
		if new_number_to_divide_left%divisible_by == 0 {
			break
		}
		new_number_to_divide_left--
	}

	right_shift_positions = new_number_to_divide_right - number_to_divide
	left_shift_positions = number_to_divide - new_number_to_divide_left

	if new_number_to_divide_right > new_number_to_divide_left {
		maximum_value = new_number_to_divide_right
	} else {
		maximum_value = new_number_to_divide_left
	}

	if right_shift_positions < left_shift_positions {
		fmt.Printf("The closest absolute number to %v and divisible by %v is %v\n", number_to_divide, divisible_by, new_number_to_divide_right)
	} else if right_shift_positions == left_shift_positions {
		fmt.Printf("Both %v and %v are the closest numbers to %v and divisible by %v\n", new_number_to_divide_right, new_number_to_divide_left, number_to_divide, divisible_by)
		fmt.Printf("However,%v is the maximum absolute value\n", maximum_value)
	} else {
		fmt.Printf("The closest absolute number to %v and divisible by %v is %v\n", number_to_divide, divisible_by, new_number_to_divide_left)
	}

}

func getAbsoluteValue(number int) int {
	if number < 0 {
		return number * -1
	}
	return number
}

func isCompletelyDivisible(number int, divisible_by int) bool {
	if number%divisible_by == 0 {
		return true
	}
	return false
}

// Find the number closest to n and divisible by m
