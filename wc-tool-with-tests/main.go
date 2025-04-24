package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var lineFlag bool
	var wordFlag bool
	var charFlag bool

	flag.BoolVar(&lineFlag, "l", false, "Count lines")
	flag.BoolVar(&wordFlag, "w", false, "Count words")
	flag.BoolVar(&charFlag, "c", false, "Count characters")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-l | -w | -c] <filename>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !lineFlag && !wordFlag && !charFlag {
		lineFlag, wordFlag, charFlag = true, true, true
	}

	for _, filePath := range flag.Args() {
		output, errMsg, exitCode := evaluateFile(filePath, lineFlag, wordFlag, charFlag)
		if exitCode != 0 {
			fmt.Fprint(os.Stderr, errMsg)
		} else {
			fmt.Print(output)
		}
	}
}

func evaluateFile(filePath string, lineFlag, wordFlag, charFlag bool) (output string, errMsg string, exitCode int) {
	file, err := os.Open(filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Sprintf("%s: %s: open: No such file or directory\n", os.Args[0], filePath), 1
		}
		if errors.Is(err, os.ErrPermission) {
			return "", fmt.Sprintf("%s: %s: open: Permission denied\n", os.Args[0], filePath), 1
		}
	}

	info, _ := file.Stat()
	if info.IsDir() {
		return "", fmt.Sprintf("%s: %s: read: Is a directory\n", os.Args[0], filePath), 1
	}

	defer file.Close()

	var result string

	if lineFlag {
		lineCount, err := lineCounter(file)
		if err != nil {
			return "", fmt.Sprintf("Error reading lines from %s: %v\n", filePath, err), 1
		}
		result += fmt.Sprintf("%8d", lineCount)
		file.Seek(0, io.SeekStart)
	}

	if wordFlag {
		wordCount, err := wordCounter(file)
		if err != nil {
			return "", fmt.Sprintf("Error reading words from %s: %v\n", filePath, err), 1
		}
		result += fmt.Sprintf("%8d", wordCount)
		file.Seek(0, io.SeekStart)
	}

	if charFlag {
		charCount, err := charCounter(file)
		if err != nil {
			return "", fmt.Sprintf("Error reading characters from %s: %v\n", filePath, err), 1
		}
		result += fmt.Sprintf("%8d", charCount)
		file.Seek(0, io.SeekStart)
	}

	result += fmt.Sprintf(" %s\n", filePath)
	return result, "", 0
}

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
