// In Go, a method is just a function with a special receiver.
// It allows functions to be associated with structs (or custom types)—similar to OOP methods in other languages.

// covering :- Basic Method on a Struct
package main

import "fmt"

type person struct {
	name    string
	age     int
	address string
	pincode int
}

func (getParameters person) printPersonDetails() {
	fmt.Printf("The name of the person is %s\n", getParameters.name)
	fmt.Printf("The age of the person is %d\n", getParameters.age)
	fmt.Printf("The address of the person is %s\n", getParameters.address)
	fmt.Printf("The pincode is %d\n", getParameters.pincode)
}
func main() {
	personObj1 := person{name: "Bob", age: 22, address: "20 Cooper Square, New York, USA", pincode: 11201}
	personObj1.printPersonDetails()
	// personObj2 := person{name: "Jhon", age: 32, address: "Metrotech Center, Brooklyn, USA", pincode: 11201}
	// personObj2.printPersonDetails()

}
