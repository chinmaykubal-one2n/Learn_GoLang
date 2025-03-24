//  A prime is a natural number (n) greater than 1 that has no positive divisors other than 1 and itself {2, 3, 5, …}

// Instead of checking till n, we can check till √n
// because a larger factor of n must be a multiple of a smaller factor that has been already checked.

package main

import (
	"fmt"
	"math"
)

func main() {
	isPrime()
}

func isPrime() {
	var number int

	fmt.Println("Enter a number:")
	fmt.Scanln(&number)

	if number <= 1 {
		fmt.Println("Please enter a valid number")
		return
	}

	for i := 2; i <= int(math.Sqrt(float64(number))); i++ {
		if number%i == 0 {
			fmt.Printf("%v is not a prime number.\n", number)
			return
		}
	}

	fmt.Printf("%v is a prime number.\n", number)

}
