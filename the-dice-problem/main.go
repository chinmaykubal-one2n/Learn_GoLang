// The idea is based on the observation that the sum of two opposite sides of a cubical dice is equal to 7.
package main

import "fmt"

func main() {
	solveDice()
}

func solveDice() {
	var dice_number int
	var opp_dice_nuber int

	fmt.Println("Enter the number on the dice:")
	fmt.Scanln(&dice_number)

	if dice_number < 1 || dice_number > 6 {
		fmt.Println("Invalid dice number. Please enter a number between 1 and 6.")
		return
	}

	opp_dice_nuber = 7 - dice_number

	fmt.Printf("The number on the opposite side of the dice is: %v\n", opp_dice_nuber)
}
