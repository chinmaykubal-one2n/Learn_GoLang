// A slice is like an array, but dynamic—it can grow and shrink in size.
// Unlike arrays, slices don’t have a fixed length.
// below stuff is also slice

package main

import "fmt"

func main() {
	numbers := [6]int{2, 3, 5, 7, 11, 13}
	var sliceOfNumbers []int = numbers[1:4]
	fmt.Println(sliceOfNumbers)

	// A slice does not store any data, it just describes a section of an underlying array.
	// Changing the elements of a slice modifies the corresponding elements of its underlying array.
	// Other slices that share the same underlying array will see those changes.
	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	fmt.Println(names)

	a := names[0:2]
	b := names[1:3]
	fmt.Println(a, b)

	b[0] = "XXX"
	fmt.Println(a, b)
	fmt.Println(names)

	// A slice literal is like an array literal without the length.
	arrayFifth := []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	fmt.Println(arrayFifth)

}
