// refacotor the below code all structs should be at first location then their corresponding objects and then functions
package main

import (
	"strings"
	"testing"
)

type lineCounterTest struct {
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

type wordCounterTest struct {
	input  string
	expect int
}

// expectations are set based on the output of wc command
var wordCounterTests = []wordCounterTest{
	{input: "Hello Golang.", expect: 2},
	{input: "This is a simple file !!", expect: 6},
	{input: "With some text inside it with some numbers 1,2,3,4,5 and some special chars.!@#$%^&*()", expect: 13},
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

type charCounterTest struct {
	input  string
	expect int
}

// expectations are set based on the output of wc command
var charCounterTests = []charCounterTest{
	{input: "Hello Golang.", expect: 13},
	{input: "This is a simple file !!", expect: 24},
	{input: "With some text inside it with some numbers 1,2,3,4,5 and some special chars.!@#$%^&*()", expect: 86},
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
