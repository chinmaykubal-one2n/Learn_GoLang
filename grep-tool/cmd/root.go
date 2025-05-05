/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var outFile string
var caseInsensitive bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "./mygrep <search_string> <filename> [-o out.txt]",
	Short: "A brief description of your application",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchString := args[0]
		var reader io.Reader

		switch len(args) {
		case 1:
			reader = os.Stdin

		case 2:
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
			err := writeStdout(matches)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
				os.Exit(1)
			}
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

	if caseInsensitive {
		searchString = strings.ToLower(searchString)
	}

	for scanner.Scan() {
		line := scanner.Text()
		compareLine := line

		if caseInsensitive {
			compareLine = strings.ToLower(line)
		}

		if strings.Contains(compareLine, searchString) {
			matches = append(matches, line)
		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func writeStdout(lines []string) error {
	for _, line := range lines {
		fmt.Println(line)
	}

	return nil
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
