package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	var lineFlag bool
	var wordFlag bool
	var charFlag bool
	var finalExitCode int

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
		file, exitCode := validateFile(filePath)
		finalExitCode = exitCode
		if file != nil {
			defer file.Close()
		}
		if file == nil {
			continue
		}
		evaluateFile(filePath, file, lineFlag, wordFlag, charFlag)
	}
	os.Exit(finalExitCode)
}

func validateFile(filePath string) (*os.File, int) {
	file, err := os.Open(filePath)
	exitCode := 0

	if errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "%s: %s: open: No such file or directory\n", os.Args[0], filePath)
		exitCode = 1
		return nil, exitCode
	}

	if errors.Is(err, os.ErrPermission) {
		fmt.Fprintf(os.Stderr, "%s: %s: open: Permission denied\n", os.Args[0], filePath)
		exitCode = 1
		return nil, exitCode
	}

	info, _ := file.Stat()
	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "%s: %s: read: Is a directory\n", os.Args[0], filePath)
		exitCode = 1
		return nil, exitCode
	}

	return file, exitCode
}

func evaluateFile(filePath string, file *os.File, lineFlag, wordFlag, charFlag bool) {
	if lineFlag {
		lineCount, lineCountErr := lineCounter(file)
		if lineCountErr != nil {
			log.Fatal(lineCountErr)
		}
		fmt.Printf("%8d", lineCount)
		file.Seek(0, io.SeekStart)
	}

	if wordFlag {
		wordCount, wordCountErr := wordCounter(file)
		if wordCountErr != nil {
			log.Fatal(wordCountErr)
		}
		fmt.Printf("%8d", wordCount)
		file.Seek(0, io.SeekStart)
	}

	if charFlag {
		charCount, charCountErr := charCounter(file)
		if charCountErr != nil {
			log.Fatal(charCountErr)
		}
		fmt.Printf("%8d", charCount)
		file.Seek(0, io.SeekStart)
	}

	fmt.Printf(" %s\n", filePath)
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
