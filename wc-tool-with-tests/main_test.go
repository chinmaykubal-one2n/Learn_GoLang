// refacotor the below code all structs should be at first location then their corresponding objects and then functions
package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

type lineCounterTest struct {
	input  string
	expect int
}

type wordCounterTest struct {
	input  string
	expect int
}

type charCounterTest struct {
	input  string
	expect int
}

var lineCounterTests = []lineCounterTest{
	{input: "", expect: 0},
	{input: "line1\n", expect: 1},
	{input: "line1\nline2\n", expect: 2},
	{input: "line1\nline2\nline3\n", expect: 3},
	{input: "line1\nline2\nline3", expect: 3},
	{input: "line1", expect: 1},
}

// expectations are set based on the output of wc command
var wordCounterTests = []wordCounterTest{
	{input: "Hello Golang.", expect: 2},
	{input: "This is a simple file !!", expect: 6},
	{input: "With some text inside it with some numbers 1,2,3,4,5 and some special chars.!@#$%^&*()", expect: 13},
}

// expectations are set based on the output of wc command
var charCounterTests = []charCounterTest{
	{input: "Hello Golang.", expect: 13},
	{input: "This is a simple file !!", expect: 24},
	{input: "With some text inside it with some numbers 1,2,3,4,5 and some special chars.!@#$%^&*()", expect: 86},
}

func TestLineCounter(t *testing.T) {
	for _, testCase := range lineCounterTests {
		reader := strings.NewReader(testCase.input)
		count, err := lineCounter(reader)
		if err != nil {
			t.Errorf("Unexpected error for input %s: %v", testCase.input, err)
		}
		if count != testCase.expect {
			t.Errorf("Expected %d lines, got %d for input %s", testCase.expect, count, testCase.input)
		}
	}
}

func TestWordCounter(t *testing.T) {
	for _, testCase := range wordCounterTests {
		reader := strings.NewReader(testCase.input)
		count, err := wordCounter(reader)
		if err != nil {
			t.Errorf("Unexpected error for input %s: %v", testCase.input, err)
		}
		if count != testCase.expect {
			t.Errorf("Expected %d words, got %d for input %s", testCase.expect, count, testCase.input)
		}
	}
}

func TestCharCounter(t *testing.T) {
	for _, testCase := range charCounterTests {
		reader := strings.NewReader(testCase.input)
		count, err := charCounter(reader)
		if err != nil {
			t.Errorf("Unexpected error for input %s: %v", testCase.input, err)
		}
		if count != testCase.expect {
			t.Errorf("Expected %d characters, got %d for input %s", testCase.expect, count, testCase.input)
		}
	}
}

func TestEvaluateFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "Hello world\nThis is Go\n"
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	allFlags := flags{
		lineFlag: true,
		wordFlag: true,
		charFlag: true,
	}

	totals := totalCounts{}

	result, errMsg, exitCode := evaluateFile(tmpFile.Name(), &allFlags, &totals)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	if errMsg != "" {
		t.Errorf("Expected empty error message, got: %s", errMsg)
	}
	if !strings.HasSuffix(result, tmpFile.Name()+"\n") {
		t.Errorf("Output did not end with filename. Got: %q", result)
	}

	expectedLineCount := 2
	expectedWordCount := 5
	expectedCharCount := len(content)

	if totals.lineCount != expectedLineCount {
		t.Errorf("Expected line count %d, got %d", expectedLineCount, totals.lineCount)
	}
	if totals.wordCount != expectedWordCount {
		t.Errorf("Expected word count %d, got %d", expectedWordCount, totals.wordCount)
	}
	if totals.charCount != expectedCharCount {
		t.Errorf("Expected char count %d, got %d", expectedCharCount, totals.charCount)
	}
}

func TestCountFromStdin(t *testing.T) {
	input := "abc\ndef ghi jkl"
	reader := bytes.NewReader([]byte(input))

	counts, err := countFromStdin(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := totalCounts{
		lineCount: 2,
		wordCount: 4,
		charCount: 15,
	}

	if counts != expected {
		t.Errorf("Expected %v, got %v", expected, counts)
	}
}
