package main

import "fmt"

func main() {

	// Accessing and Updating Values
	myMap := make(map[string]int)
	myMap["age"] = 20
	myMap["age"] = 30
	fmt.Println(myMap)

	myMap2 := map[string]int{"apple": 100, "mango": 200}
	fmt.Println(myMap2["apple"], myMap2["mango"])

	//  Checking if a Key Exists
	price, isItThere := myMap2["banana"]

	if isItThere {
		fmt.Println("The price is ", price)
	} else {
		fmt.Println("fruit not found")
	}

	// Deleting a Key from a Map
	fmt.Println("before deleting the apple", myMap2)
	delete(myMap2, "apple")
	delete(myMap2, "CHEERY") // If the key doesn’t exist, delete() does nothing (no error).
	fmt.Println("after deleting the apple", myMap2)

	// Iterating Over a Map
	myMap3 := map[string]int{"apple": 100, "banana": 50, "omlete": 30, "vadapav": 20}

	for key, value := range myMap3 {
		fmt.Println(key, "-->", value)
	}

}
