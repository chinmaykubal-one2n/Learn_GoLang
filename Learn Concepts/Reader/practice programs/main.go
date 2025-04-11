// Open a file
// Read its contents using io.Reader
// Count how many words are in it
// Handle all possible errors
// Use defer to safely close the file

// Requirements:
// Use os.Open() to open the file
// Use defer file.Close()
// Use a Reader to read file data
// Count how many words are in the file
// Handle errors gracefully
// Print the word count at the end

package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var filePath string = "./content.txt"
	file, error := os.Open(filePath)

	if error != nil {
		fmt.Println("Error opening the file:", error)
		os.Exit(1)
	}
	defer file.Close()
	printFile(file)
	os.Exit(0)
}

func printFile(file io.Reader) {
	buffer := make([]byte, 4096)

	for {
		noOfBytes, error := file.Read(buffer)

		if error == io.EOF {
			break
		}

		fmt.Print(string(buffer[:noOfBytes]))
	}
}
