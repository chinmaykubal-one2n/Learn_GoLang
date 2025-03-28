// return only the prime numbers from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
// Sample Output: 2, 3, 5, 7
package main

import (
	"fmt"
	"math"
)

func isPrime(number int) bool {
	if number < 2 {
		return false
	}

	for i := 2; i <= int(math.Sqrt(float64(number))); i++ {
		if number%i == 0 {
			return false
		}
	}

	return true
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var primes []int
	for _, number := range numbers {
		if isPrime(number) {
			primes = append(primes, number)
		}
	}
	fmt.Println("Sample Output:", primes)
}
