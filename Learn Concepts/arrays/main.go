// An array in Golang is a fixed-size collection of elements of the same type.

package main

import "fmt"

func main() {

	var arrayFirst []int                                                               // slice
	var arraySecond [5]int                                                             // array
	arrayThird := [5]int{1, 2, 3, 4, 5}                                                // array
	arrayFourth := [...]int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20} // array
	arrayFifth := []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}     // working is same as above // slice
	fmt.Println(arrayFirst, arraySecond, arrayThird, arrayFourth, arrayFifth)

	arraySixth := []string{"a", "b", "c", "d", "e"} // slice
	arraySixth[1] = "z"
	fmt.Printf("arraySixth: %v\n", arraySixth)

	arraySeventh := [3]int{99, 98, 97} // array
	for i := 0; i < len(arraySeventh); i++ {
		fmt.Println(arraySeventh[i])
	}

	for index, value := range arraySeventh {
		fmt.Println("Index:", index, "Value:", value)
	}

}
