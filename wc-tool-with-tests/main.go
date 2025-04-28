package main

import (
	"bufio"
	"bytes"
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
		readFromStdin(&allFlags)
		os.Exit(osExitCode)
		// fmt.Fprintf(os.Stderr, "Usage: %s [-l | -w | -c] <filename>\n", os.Args[0])
		// flag.PrintDefaults()
		// os.Exit(1)
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
		printAllCounts(&allFlags, &totals)
		fmt.Println(" total")
	}

	os.Exit(osExitCode)
}

func printAllCounts(allFlags *flags, totalCounts *totalCounts) {
	if allFlags.lineFlag {
		fmt.Printf("%8d", totalCounts.lineCount)
	}
	if allFlags.wordFlag {
		fmt.Printf("%8d", totalCounts.wordCount)
	}
	if allFlags.charFlag {
		fmt.Printf("%8d", totalCounts.charCount)
	}
}

func evaluateFile(filePath string, allFlags *flags, totals *totalCounts) (output string, errMsg string, exitCode int) {
	const errorCode = 1
	const successCode = 0
	const emptyString = ""

	file, errMsg, exitCode := validateFile(filePath)
	if exitCode != successCode {
		return emptyString, errMsg, exitCode
	}
	defer file.Close()

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

func validateFile(filePath string) (*os.File, string, int) {
	const errorCode = 1
	const successCode = 0
	const emptyString = ""

	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Sprintf("%s: %s: open: No such file or directory\n", os.Args[0], filePath), errorCode
		}
		if errors.Is(err, os.ErrPermission) {
			return nil, fmt.Sprintf("%s: %s: open: Permission denied\n", os.Args[0], filePath), errorCode
		}
	}

	info, _ := file.Stat()
	if info.IsDir() {
		file.Close()
		return nil, fmt.Sprintf("%s: %s: read: Is a directory\n", os.Args[0], filePath), errorCode
	}

	return file, emptyString, successCode
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

func readFromStdin(allFlags *flags) {
	totalCounts, err := countFromStdin(os.Stdin)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if !allFlags.lineFlag && !allFlags.wordFlag && !allFlags.charFlag {
		allFlags.lineFlag, allFlags.wordFlag, allFlags.charFlag = true, true, true
	}

	fmt.Println()
	printAllCounts(allFlags, &totalCounts)
	fmt.Println()
}

func countFromStdin(stdInput io.Reader) (totalCounts, error) {
	input, err := io.ReadAll(stdInput)
	if err != nil {
		return totalCounts{}, err
	}

	lines, err := lineCounter(bytes.NewReader(input))
	if err != nil {
		return totalCounts{}, err
	}

	words, err := wordCounter(bytes.NewReader(input))
	if err != nil {
		return totalCounts{}, err
	}

	chars, err := charCounter(bytes.NewReader(input))
	if err != nil {
		return totalCounts{}, err
	}

	return totalCounts{
		lineCount: lines,
		wordCount: words,
		charCount: chars,
	}, nil
}
