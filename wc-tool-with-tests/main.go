package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
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

var mu sync.Mutex

func main() {
	var osExitCode int
	allFlags := flags{}
	totals := totalCounts{}

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

	var wg sync.WaitGroup

	for _, filePath := range flag.Args() {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			output, errMsg, exitCode := evaluateFile(path, &allFlags, &totals)
			if exitCode != 0 {
				mu.Lock()
				fmt.Fprint(os.Stderr, errMsg)
				osExitCode = exitCode
				mu.Unlock()
			} else {
				mu.Lock()
				fmt.Print(output)
				mu.Unlock()
			}
		}(filePath)
	}

	wg.Wait()

	if flag.NArg() > 1 {
		if allFlags.lineFlag {
			fmt.Printf("%8d", totals.lineCount)
		}
		if allFlags.wordFlag {
			fmt.Printf("%8d", totals.wordCount)
		}
		if allFlags.charFlag {
			fmt.Printf("%8d", totals.charCount)
		}
		fmt.Println(" total")
	}

	os.Exit(osExitCode)
}

func evaluateFile(filePath string, allFlags *flags, totals *totalCounts) (output string, errMsg string, exitCode int) {
	const errorCode = 1
	const successCode = 0
	const emptyString = ""

	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return emptyString, fmt.Sprintf("%s: %s: open: No such file or directory\n", os.Args[0], filePath), errorCode
		}
		if errors.Is(err, os.ErrPermission) {
			return emptyString, fmt.Sprintf("%s: %s: open: Permission denied\n", os.Args[0], filePath), errorCode
		}
		return emptyString, fmt.Sprintf("%s: %s: open: %v\n", os.Args[0], filePath, err), errorCode
	}
	defer file.Close()

	info, _ := file.Stat()
	if info.IsDir() {
		return emptyString, fmt.Sprintf("%s: %s: read: Is a directory\n", os.Args[0], filePath), errorCode
	}

	var result string

	if allFlags.lineFlag {
		lineCount, err := lineCounter(file)
		if err != nil {
			return emptyString, fmt.Sprintf("Error reading lines from %s: %v\n", filePath, err), errorCode
		}
		mu.Lock()
		totals.lineCount += lineCount
		mu.Unlock()
		result += fmt.Sprintf("%8d", lineCount)
		file.Seek(0, io.SeekStart)
	}

	if allFlags.wordFlag {
		wordCount, err := wordCounter(file)
		if err != nil {
			return emptyString, fmt.Sprintf("Error reading words from %s: %v\n", filePath, err), errorCode
		}
		mu.Lock()
		totals.wordCount += wordCount
		mu.Unlock()
		result += fmt.Sprintf("%8d", wordCount)
		file.Seek(0, io.SeekStart)
	}

	if allFlags.charFlag {
		charCount, err := charCounter(file)
		if err != nil {
			return emptyString, fmt.Sprintf("Error reading characters from %s: %v\n", filePath, err), errorCode
		}
		mu.Lock()
		totals.charCount += charCount
		mu.Unlock()
		result += fmt.Sprintf("%8d", charCount)
	}

	result += fmt.Sprintf(" %s\n", filePath)
	return result, emptyString, successCode
}

func genericCounter(file io.Reader, split bufio.SplitFunc) (int, error) {
	count := 0
	scanner := bufio.NewScanner(file)
	scanner.Split(split)

	for scanner.Scan() {
		count++
	}

	err := scanner.Err()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func lineCounter(file io.Reader) (int, error) {
	return genericCounter(file, bufio.ScanLines)
}

func wordCounter(file io.Reader) (int, error) {
	return genericCounter(file, bufio.ScanWords)
}

func charCounter(file io.Reader) (int, error) {
	return genericCounter(file, bufio.ScanRunes)
}
