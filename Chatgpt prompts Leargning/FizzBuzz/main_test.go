// implement a function FizzBuzz(n int) []string that returns a slice of strings from 1 to n, where:
// For multiples of 3, add "Fizz" instead of the number
// For multiples of 5, add "Buzz"
// For numbers which are multiples of both 3 and 5, add "FizzBuzz"
// Otherwise, return the number as a string
// Example :- FizzBuzz(5) == []string{"1", "2", "Fizz", "4", "Buzz"}

package main

import (
	"reflect"
	"testing"
)

type testCaseType struct {
	input  int
	output []string
}

var testCases = []testCaseType{
	{
		input:  5,
		output: []string{"1", "2", "Fizz", "4", "Buzz"},
	},
	{
		input:  10,
		output: []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz"},
	},
	{
		input:  15,
		output: []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"},
	},
	{
		input:  0,
		output: []string{},
	},
}

func TestFizzBuzz(t *testing.T) {
	for _, testCase := range testCases {
		output := FizzBuzz(testCase.input)
		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("For Input %d Expected %#v, but got %#v", testCase.input, testCase.output, output)
		}
	}
}
