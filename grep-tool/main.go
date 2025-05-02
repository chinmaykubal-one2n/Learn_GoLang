package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	const (
		successCode int = 0
		errorCode   int = 1
	)

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <search_string> <filename>\n", os.Args[0])
		os.Exit(errorCode)
	}

	searchString := args[0]

	if len(args) == 1 {
		lines, err := grepReader(searchString, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(errorCode)
		}
		fmt.Println()
		fmt.Println()
		for _, line := range lines {
			fmt.Printf("%s\n", line)
		}
	}

	if len(args) == 2 {
		filename := args[1]
		file, err := validateFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(errorCode)
		}
		defer file.Close()

		matches, err := grepReader(searchString, file)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %v\n", os.Args[0], filename, err)
			os.Exit(errorCode)
		}

		for _, matchedLines := range matches {
			fmt.Println(matchedLines)
		}
	}

	os.Exit(successCode)
}

func validateFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)

	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%s: %s: Permission denied", os.Args[0], filename)
		}
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: %s: open: No such file or directory", os.Args[0], filename)
		}
		return nil, fmt.Errorf("%s: %s: %v", os.Args[0], filename, err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("%s: %s: %v", os.Args[0], filename, err)
	}

	if info.IsDir() {
		file.Close()
		return nil, fmt.Errorf("%s: %s: read: Is a directory", os.Args[0], filename)
	}

	return file, nil
}

func grepReader(searchString string, reader io.Reader) ([]string, error) {
	var matches []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchString) {
			matches = append(matches, line)
		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return matches, nil
}
