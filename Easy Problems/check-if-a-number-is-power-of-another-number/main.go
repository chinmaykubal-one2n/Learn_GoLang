// y=x^n --> using log on both sides...(n is an integer) --> log(y)=log(x^n) --> log(y)=nlog(x)
// 8=2^n  .. where n=3 which our program will find
// n = log(y)/log(x) now : We compute the ratio n = log(y)/log(x)​.
// And if this ratio is very close to an integer (within a small tolerance like 1e-10), then y is a some nth power of x.
package main

import (
	"fmt"
	"math"
)

var x, y int

func main() {
	if isPower() {
		fmt.Printf("Yes, %v is a power of %v\n", y, x)
	} else {
		fmt.Printf("No, %v is not a power of %v\n", y, x)
	}
}

func isPower() bool {
	fmt.Println("Enter x: ")
	fmt.Scan(&x)
	fmt.Println("Enter y: ")
	fmt.Scan(&y)

	if x <= 0 || y <= 0 {
		return false
	}

	// Below are the maths formulae to check if y is a power of x
	n := math.Log(float64(y)) / math.Log(float64(x))
	return math.Abs(n-math.Round(n)) < 1e-10
}
