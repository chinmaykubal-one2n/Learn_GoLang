// Generics in Go allow functions and types to work with any data type
// instead of being limited to a specific type like int or string.
// This makes code more reusable and type-safe. Introduced in Go 1.18

package main

import "fmt"

func print[T any](value T) {
	fmt.Printf("value is %v\n", value)
}

func Add[T int | float64](num1, num2 T) T {
	return num1 + num2
}

func main() {
	print("Hello")
	print(true)
	print(90)
	fmt.Printf("Addition is %v\n", Add(3, 5.0))
}
