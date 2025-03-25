package main

import "fmt"

func main() {

	var variable1 int = 10

	variable2 := &variable1 // output will be the address of variable1

	fmt.Println("Value of variable1 is: ", variable1)    // output will be 10
	fmt.Println("Address of variable1 is: ", &variable1) // output will be the address of variable1
	fmt.Println("Value of variable2 is: ", variable2)    // output will be the address of variable1
	fmt.Println("Value of variable2 is: ", *variable2)   // output will be 10, value at the address of variable1
	fmt.Println("Changing value of variable1 using variable2")

	*variable2 = 20 // Dereferencing the pointer to change the value at the address of variable1

	fmt.Println("Value of variable1 is: ", variable1)  // output will be 20
	fmt.Println("Value of variable2 is: ", *variable2) // output will be 20
}

// & (Address-of) → Gets the memory address of a variable.
// * (Dereference) → Access/modify the value stored at that address.
