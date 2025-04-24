package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type totalCounts struct {
	lineCount int
	wordCount int
	charCount int
}

type flags struct {
	lineFlag bool
	wordFlag bool
	charFlag bool
}

func main() {
	var osExitCode int
	allFlags := flags{}          // all flags are false by default
	totalCounts := totalCounts{} // all counts are 0 by default

	flag.BoolVar(&allFlags.lineFlag, "l", false, "Count lines")
	flag.BoolVar(&allFlags.wordFlag, "w", false, "Count words")
	flag.BoolVar(&allFlags.charFlag, "c", false, "Count characters")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-l | -w | -c] <filename>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !allFlags.lineFlag && !allFlags.wordFlag && !allFlags.charFlag {
		allFlags.lineFlag, allFlags.wordFlag, allFlags.charFlag = true, true, true
	}

	for _, filePath := range flag.Args() {
		output, errMsg, exitCode := evaluateFile(filePath, &allFlags, &totalCounts)
		if exitCode != 0 {
			osExitCode = exitCode
			fmt.Fprint(os.Stderr, errMsg)
		} else {
			fmt.Print(output)
		}
	}

	if flag.NArg() > 1 {
		if allFlags.lineFlag {
			fmt.Printf("%8d", totalCounts.lineCount)
		}
		if allFlags.wordFlag {
			fmt.Printf("%8d", totalCounts.wordCount)
		}
		if allFlags.charFlag {
			fmt.Printf("%8d", totalCounts.charCount)
		}
		fmt.Printf(" total\n")
	}

	os.Exit(osExitCode)
}

func evaluateFile(filePath string, allFlags *flags, totalCounts *totalCounts) (output string, errMsg string, exitCode int) {
	const errorCode = 1
	const successCode = 0

	file, err := os.Open(filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Sprintf("%s: %s: open: No such file or directory\n", os.Args[0], filePath), errorCode
		}
		if errors.Is(err, os.ErrPermission) {
			return "", fmt.Sprintf("%s: %s: open: Permission denied\n", os.Args[0], filePath), errorCode
		}
	}

	info, _ := file.Stat()
	if info.IsDir() {
		return "", fmt.Sprintf("%s: %s: read: Is a directory\n", os.Args[0], filePath), errorCode
	}

	defer file.Close()

	var result string

	if allFlags.lineFlag {
		lineCount, err := lineCounter(file)
		if err != nil {
			return "", fmt.Sprintf("Error reading lines from %s: %v\n", filePath, err), errorCode
		}
		totalCounts.lineCount += lineCount
		result += fmt.Sprintf("%8d", lineCount)
		file.Seek(0, io.SeekStart)
	}

	if allFlags.wordFlag {
		wordCount, err := wordCounter(file)
		if err != nil {
			return "", fmt.Sprintf("Error reading words from %s: %v\n", filePath, err), errorCode
		}
		totalCounts.wordCount += wordCount
		result += fmt.Sprintf("%8d", wordCount)
		file.Seek(0, io.SeekStart)
	}

	if allFlags.charFlag {
		charCount, err := charCounter(file)
		if err != nil {
			return "", fmt.Sprintf("Error reading characters from %s: %v\n", filePath, err), errorCode
		}
		totalCounts.charCount += charCount
		result += fmt.Sprintf("%8d", charCount)
		file.Seek(0, io.SeekStart)
	}

	result += fmt.Sprintf(" %s\n", filePath)
	return result, "", successCode
}

// combine the three functions into one
func lineCounter(file io.Reader) (int, error) {
	lineCount := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lineCount++
	}

	err := scanner.Err()
	if err != nil {
		return 0, err
	}

	return lineCount, nil
}

func wordCounter(file io.Reader) (int, error) {
	wordCount := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wordCount++
	}

	err := scanner.Err()
	if err != nil {
		return 0, err
	}

	return wordCount, nil
}

func charCounter(file io.Reader) (int, error) {
	charCount := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		charCount++
	}

	err := scanner.Err()
	if err != nil {
		return 0, err
	}

	return charCount, nil
}
