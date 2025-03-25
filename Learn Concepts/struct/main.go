package main

import "fmt"

type Person struct {
	Name    string
	Age     int
	Address string
	Salary  int
}

func main() {
	personDetails := Person{"Bob", 25, "New York", 90000}
	personDetails2 := Person{Name: "Tom", Age: 35, Address: "North Carolina"} // Most readable way and if didnt write the field name then it will be zero value
	personDetails3 := &Person{"BigBoy", 45, "South Carolina", 400000}

	fmt.Println(personDetails)   //output is {Bob 25 New York 90000}
	fmt.Println(personDetails2)  //output is {Tom 35 North Carolina 0}
	fmt.Println(personDetails3)  //output is &{BigBoy 45 South Carolina 400000}
	fmt.Println(*personDetails3) //output is {BigBoy 45 South Carolina 400000}

	// Accessing struct fields
	// fmt.Printf("Name: %s\n", personDetails.Name)
	// fmt.Printf("Age: %d\n", personDetails.Age)
	// fmt.Printf("Address: %s\n", personDetails.Address)
	// fmt.Printf("Salary: %d\n", personDetails.Salary)

	// Updating struct fields via pointer
	// so technically we have changed the values of the personDetails variable itself
	updatePersonDetails(&personDetails)
	fmt.Println(personDetails)  //output is {Alice 30 California 100000}
	fmt.Println(personDetails2) //output is {Tom 35 North Carolina 200000}

	// fmt.Printf("Name: %s\n", personDetails.Name)
	// fmt.Printf("Age: %d\n", personDetails.Age)
	// fmt.Printf("Address: %s\n", personDetails.Address)
	// fmt.Printf("Salary: %d\n", personDetails.Salary)
}

func updatePersonDetails(passVariableHere *Person) {
	passVariableHere.Name = "Alice"
	passVariableHere.Age = 30
	passVariableHere.Address = "California"
	passVariableHere.Salary = 100000
}
