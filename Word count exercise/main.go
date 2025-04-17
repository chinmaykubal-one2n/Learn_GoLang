package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	filePath := "example_file.txt"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lineCount, err := lineCounter(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total lines: %d\n", lineCount)
}

func lineCounter(file io.Reader) (int, error) {
	buffer := make([]byte, 32*1024)
	lineCount := 0
	lineSeperator := []byte{'\n'}

	for {
		noOfBytes, err := file.Read(buffer)
		lineCount += bytes.Count(buffer[:noOfBytes], lineSeperator)

		switch {
		case err == io.EOF:
			return lineCount, nil

		case err != nil:
			return lineCount, err
		}
	}

}

// after writing raise a pr for the same, lets write some more code then raise the pr.
