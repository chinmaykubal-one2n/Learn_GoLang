package main

import (
	"log"
	"os"
)

func main() {
	file := "./example_file.txt"
	readFile(file)
}

func readFile(file string) {
	data, error := os.ReadFile(file)

	if error != nil {
		log.Fatal(error)
	}

	os.Stdout.Write(data)
}

// after writing raise a pr for the same.
