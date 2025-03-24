package main

import "fmt"

func main() {
	sum := sumOfFirstNNaturalNumbers()
	fmt.Println(sum)
}

func sumOfFirstNNaturalNumbers() int {
	var number int
	sum := 0

	fmt.Println("Enter a number:")
	fmt.Scanln(&number)

	for i := 1; i <= number; i++ {
		sum = sum + i
	}

	return sum

}
