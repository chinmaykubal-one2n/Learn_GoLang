package main

import (
	"fmt"
	"io"
	"os"
)

func printData(r io.Reader) {
	buf := make([]byte, 4)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Print(string(buf[:n]))
	}
}

func main() {
	// printData(strings.NewReader("hello")) // from string
	file, err := os.Open("./lorem_ipsum.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	printData(file) // from file

}
