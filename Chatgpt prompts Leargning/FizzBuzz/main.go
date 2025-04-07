// implement a function FizzBuzz(n int) []string that returns a slice of strings from 1 to n, where:
// For multiples of 3, add "Fizz" instead of the number
// For multiples of 5, add "Buzz"
// For numbers which are multiples of both 3 and 5, add "FizzBuzz"
// Otherwise, return the number as a string
// Example :- FizzBuzz(5) == []string{"1", "2", "Fizz", "4", "Buzz"}

package main

import (
	"fmt"
	"strconv"
)

func FizzBuzz(n int) []string {
	var output []string
	const (
		fizz     = "Fizz"
		buzz     = "Buzz"
		fizzbuzz = "FizzBuzz"
	)

	if n <= 0 {
		return []string{}
	}

	for i := 1; i <= n; i++ {
		if i%3 == 0 && i%5 == 0 {
			output = append(output, fizzbuzz)
		} else if i%3 == 0 {
			output = append(output, fizz)
		} else if i%5 == 0 {
			output = append(output, buzz)
		} else {
			output = append(output, strconv.Itoa(i))
		}
	}

	return output
}

func main() {
	number := 20
	output := FizzBuzz(number)
	fmt.Printf("output is %#v\n", output)
}
