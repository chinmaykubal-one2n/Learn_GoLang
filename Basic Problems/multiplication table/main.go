package main

import "fmt"

func main() {
	tableOfNumber := multiplicationTable()
	fmt.Println(tableOfNumber)

	maultiplicationTablePrint()
}

func multiplicationTable() []int {
	var number int
	var tableOfNumber []int

	fmt.Println("Enter a number:")
	fmt.Scanln(&number)

	for i := 1; i <= 10; i++ {
		calculation := number * i
		tableOfNumber = append(tableOfNumber, calculation)
	}

	return tableOfNumber
}

func maultiplicationTablePrint() {
	var number int

	fmt.Println("Enter a number:")
	fmt.Scanln(&number)

	for i := 1; i <= 10; i++ {
		fmt.Printf("%v x %v = %v\n", number, i, number*i)
	}
}
