package main

import (
	"fmt"
)

type Transport interface {
	BookTicket(name string) string
	GetPrice() float64
}

type Bus struct {
	Route  string
	Booked bool
}

func (b *Bus) BookTicket(name string) string {
	b.Booked = true
	return fmt.Sprintf("Ticket booked for %s on Bus (%s route)", name, b.Route)
}

func (b *Bus) GetPrice() float64 {
	return 100.0
}

type Train struct {
	Coach  string
	Booked bool
}

func (t *Train) BookTicket(name string) string {
	t.Booked = true
	return fmt.Sprintf("Ticket booked for %s on Train (Coach %s)", name, t.Coach)
}

func (t *Train) GetPrice() float64 {
	return 250.0
}

type Flight struct {
	FlightNumber string
	Booked       bool
}

func (f *Flight) BookTicket(name string) string {
	f.Booked = true
	return fmt.Sprintf("Ticket booked for %s on Flight (%s)", name, f.FlightNumber)
}

func (f *Flight) GetPrice() float64 {
	return 500.0
}

func ProcessBooking(t Transport, name string) {
	fmt.Println(t.BookTicket(name))
	fmt.Printf("Total cost: ₹%.2f\n", t.GetPrice())
}

func main() {
	bus := &Bus{Route: "City Center"}
	train := &Train{Coach: "S2"}
	flight := &Flight{FlightNumber: "AI202"}

	ProcessBooking(bus, "Alice")
	ProcessBooking(train, "Bob")
	ProcessBooking(flight, "Charlie")
}
