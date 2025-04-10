// Challenge: Shape Interface
// You’ll define an interface Shape with two methods, Area(), Perimeter()
// then create two shapes (Circle and Rectangle) that implement this interface.

package main

import "fmt"

const pi = 3.14

type shape interface {
	area() float64
	perimeter() float64
}

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return pi * c.radius * c.radius
}

func (r circle) perimeter() float64 {
	return 2 * pi * r.radius
}

type rectangle struct {
	width  float64
	length float64
}

func (r rectangle) area() float64 {
	return r.length * r.width
}

func (r rectangle) perimeter() float64 {
	return 2 * (r.length + r.width)
}

func calcArea(s shape) float64 {
	return s.area()
}

func calcPerimeter(s shape) float64 {
	return s.perimeter()
}

func main() {
	circleObj := circle{
		radius: 88,
	}

	fmt.Printf("Area of circle is %.2f\n", calcArea(circleObj))
	fmt.Printf("Perimeter of circle is %.2f\n", calcPerimeter(circleObj))

	rectangleObj := rectangle{
		width:  20,
		length: 30,
	}

	fmt.Printf("Area of rectangle is %.2f\n", calcArea(rectangleObj))
	fmt.Printf("Perimeter of rectangle is %.2f\n", calcPerimeter(rectangleObj))
}
