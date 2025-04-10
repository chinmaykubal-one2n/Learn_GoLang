// An interface in Go is a type that defines a set of methods.
// Any type that has those methods automatically "implements" that interface — no keyword needed!

package main

import "fmt"

type Database interface {
	Save(data string)
}

type MySQL struct{}

func (m MySQL) Save(data string) {
	fmt.Println("Saved to MySQL:", data)
}

type MongoDB struct{}

func (m MongoDB) Save(data string) {
	fmt.Println("Saved to MongoDB:", data)
}

func Store(d Database, info string) {
	d.Save(info)
}

func main() {
	Store(MySQL{}, "User123")
	Store(MongoDB{}, "Post456")
}
