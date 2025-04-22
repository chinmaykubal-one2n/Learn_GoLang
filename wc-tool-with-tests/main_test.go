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

var wordCounterTests = []wordCounterTest{
	{input: "", expect: 0},
	{input: "hello", expect: 1},
	{input: "hello world", expect: 2},
	{input: "hello   world", expect: 2},
	{input: "hello\nworld", expect: 2},
	{input: "hello\tworld\nfoo", expect: 3},
	{input: "hello\tworld\nfoo\n", expect: 3},
	{input: "   spaced   out    words   ", expect: 3},
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
