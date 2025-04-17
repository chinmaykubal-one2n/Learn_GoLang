/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	lineFlag bool
	wordFlag bool
	charFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wc [flags] <filename>",
	Short: "A mini version of the wc command",
	Long:  `Trying to mimic the command line program that implements Unix wc like functionality`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exitCode := 0
		totalLineCount := 0
		totalWordCount := 0
		totalCharCount := 0

		for _, filePath := range args {
			file, err := os.Open(filePath)

			if err != nil && errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "%s: %s: open: No such file or directory\n", os.Args[0], filePath)
				exitCode++
				continue
			}

			if err != nil && errors.Is(err, os.ErrPermission) {
				fmt.Fprintf(os.Stderr, "%s: %s: open: Permission denied\n", os.Args[0], filePath)
				exitCode++
				continue
			}

			info, _ := file.Stat()
			if info.IsDir() {
				fmt.Fprintf(os.Stderr, "%s: %s: read: Is a directory\n", os.Args[0], filePath)
				exitCode++
				continue
			}

			defer file.Close()

			if !lineFlag && !wordFlag && !charFlag {
				lineFlag = true
				wordFlag = true
				charFlag = true
			}

			if lineFlag {
				lineCount, err := lineCounter(file)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%8d", lineCount)
				totalLineCount = totalLineCount + lineCount
				file.Seek(0, io.SeekStart)
			}

			if wordFlag {
				wordCount, err := wordCounter(file)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%8d", wordCount)
				totalWordCount = totalWordCount + wordCount
				file.Seek(0, io.SeekStart)
			}

			if charFlag {
				charCount, err := charCounter(file)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%8d", charCount)
				totalCharCount = totalCharCount + charCount
			}
			fmt.Printf(" %s\n", filePath)

		}

		if lineFlag && totalLineCount > 0 {
			fmt.Printf("%8d", totalLineCount)
		}
		if wordFlag && totalWordCount > 0 {
			fmt.Printf("%8d", totalWordCount)
		}
		if charFlag && totalCharCount > 0 {
			fmt.Printf("%8d", totalCharCount)
		}
		fmt.Printf(" total\n")

		os.Exit(exitCode)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wc-tool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolVarP(&lineFlag, "lines", "l", false, "Count lines in the file (like wc -l)")
	rootCmd.Flags().BoolVarP(&wordFlag, "words", "w", false, "Count words in the file (like wc -w)")
	rootCmd.Flags().BoolVarP(&charFlag, "characters", "c", false, "Count characters in the file (like wc -c)")
}

// old approach
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

func wordCounter(file io.Reader) (int, error) {
	wordCount := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) //bufio.ScanLines (splits on a new line, we can make use this to calc. no. of lines)

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
