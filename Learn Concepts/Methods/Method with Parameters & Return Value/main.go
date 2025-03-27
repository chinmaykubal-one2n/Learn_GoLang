// Methods can take parameters and return values, just like regular functions.
package main

import "fmt"

type rectangle struct {
	width, height float64
}

func (getParameters rectangle) calcAreaOfRectangle() float64 {
	return getParameters.height * getParameters.width
}

func main() {
	objOfRectangle := rectangle{width: 10, height: 20}
	areaOfRectangle := objOfRectangle.calcAreaOfRectangle()
	fmt.Printf("Area of Rectangle is %v \n", areaOfRectangle)
}
