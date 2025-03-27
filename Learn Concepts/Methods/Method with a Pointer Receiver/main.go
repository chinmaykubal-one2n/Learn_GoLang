// Method with a Pointer Receiver
// Use a pointer receiver (*counter) to modify the original struct.
// Without *counter, the struct won’t change (Go passes a copy).
package main

import "fmt"

type counter struct {
	value int
}

// Method with a pointer receiver (`*counter`) to modify the struct
func (getValue *counter) increment() {
	getValue.value++
}

func main() {
	objOfCounter := counter{value: 10}
	objOfCounter.increment()
	fmt.Printf("incremented value is %d\n", objOfCounter.value)
}
