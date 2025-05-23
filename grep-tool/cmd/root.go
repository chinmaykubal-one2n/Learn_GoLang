/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var outFile string
var caseInsensitive bool
var recursive bool
var after, before int
var countOnly bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "./mygrep <search_string> <filename> [-o out.txt]",
	Short: "A grep-like command-line tool written in Go",
	Long: `mygrep allows searching for text in files or directories, 
	with options like case-insensitive search and output redirection.`,
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		searchString := args[0]
		var reader io.Reader

		if recursive {
			var filename string
			if len(args) == 1 {
				filename = "."
			} else {
				filename = args[1]
			}
			recursiveSearch(searchString, filename, os.Stdout)
			return
		}

		if len(args) > 2 {
			for _, arg := range args {
				if arg == searchString {
					continue
				}
				filename := arg
				file, err := validateFile(filename)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				defer file.Close()
				reader = file

				matches, err := grepReader(searchString, reader)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
					os.Exit(1)
				}
				writeStdout(matches, os.Stdout)
			}
			return
		}

		if len(args) == 1 {
			reader = os.Stdin
		}

		if len(args) == 2 {
			filename := args[1]
			file, err := validateFile(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer file.Close()
			reader = file
		}

		matches, err := grepReader(searchString, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
			os.Exit(1)
		}

		if outFile == "" {
			writeStdout(matches, os.Stdout)
		} else {
			err := writeToFile(outFile, matches)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&outFile, "out", "o", "", "Write output to file instead of stdout")
	rootCmd.Flags().BoolVarP(&caseInsensitive, "i", "i", false, "Ignore case when searching")
	rootCmd.Flags().BoolVarP(&recursive, "r", "r", false, "Search recursively in directories")
	rootCmd.Flags().IntVarP(&after, "after", "A", 0, "Print n lines after match")
	rootCmd.Flags().IntVarP(&before, "before", "B", 0, "Print n lines before match")
	rootCmd.Flags().BoolVarP(&countOnly, "count", "c", false, "Only print count of matches")
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
	}

	info, _ := file.Stat()

	if info.IsDir() {
		file.Close()
		return nil, fmt.Errorf("%s: %s: read: Is a directory", os.Args[0], filename)
	}

	return file, nil
}

func grepReader(searchString string, reader io.Reader) ([]string, error) {
	var matches []string
	var count int

	scanner := bufio.NewScanner(reader)

	if caseInsensitive {
		searchString = strings.ToLower(searchString)
	}

	beforeBuffer := make([]string, 0, before)
	afterRemaining := 0

	for scanner.Scan() {
		line := scanner.Text()
		compareLine := line
		if caseInsensitive {
			compareLine = strings.ToLower(line)
		}

		match := strings.Contains(compareLine, searchString)

		if match {
			count++
			if countOnly {
				continue
			}

			for _, b := range beforeBuffer {
				matches = append(matches, b)
			}

			beforeBuffer = beforeBuffer[:0]

			matches = append(matches, line)

			afterRemaining = after

		} else if afterRemaining > 0 {
			matches = append(matches, line)
			afterRemaining--

		} else if before > 0 {
			if len(beforeBuffer) == before {
				beforeBuffer = beforeBuffer[1:]
			}
			beforeBuffer = append(beforeBuffer, line)
		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	if countOnly {
		return []string{strconv.Itoa(count)}, nil
	}

	return matches, nil
}

func writeStdout(lines []string, out io.Writer) {
	for _, line := range lines {
		// fmt.Println(line)
		fmt.Fprintf(out, "%s\n", line)
	}
}

func writeToFile(outPath string, lines []string) error {
	_, err := os.Stat(outPath)
	if err == nil {
		return fmt.Errorf("%s already exists", outPath)
	}

	file, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteToStdOutForRecursiveFiles(w io.Writer, lines []string, filename string) {
	for _, line := range lines {
		fmt.Fprintf(w, "%s:%s\n", filename, line)
	}
}

func recursiveSearch(searchString, root string, out io.Writer) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			file, err := validateFile(path)
			if err != nil {
				return
			}
			defer file.Close()

			matches, err := grepReader(searchString, file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
				return
			}

			if len(matches) > 0 {
				mu.Lock()
				WriteToStdOutForRecursiveFiles(out, matches, path)
				mu.Unlock()
			}
		}(path)

		return nil
	})

	wg.Wait()
}
