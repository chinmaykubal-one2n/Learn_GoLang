// https://www.geeksforgeeks.org/program-calculate-distance-two-points/

package main

import (
	"fmt"
	"math"
)

func main() {
	distance()
}

func distance() {
	var x1, y1, x2, y2 float64
	var distance float64

	fmt.Println("Enter coordinates:")
	fmt.Println("Enter x1: ")
	fmt.Scan(&x1)
	fmt.Println("Enter y1: ")
	fmt.Scan(&y1)
	fmt.Println("Enter x2: ")
	fmt.Scan(&x2)
	fmt.Println("Enter y2: ")
	fmt.Scan(&y2)

	distance = math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2)) // distance formula by pythagoras theorem

	fmt.Printf("Distance between (%v, %v) and (%v, %v) is %v\n", x1, y1, x2, y2, distance)
}
