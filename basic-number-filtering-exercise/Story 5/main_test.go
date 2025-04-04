// return only the even and multiples of 5 from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20
// Sample Output: 10, 20

package main

import (
	"reflect"
	"testing"
)

type testCaseType struct {
	input  []int
	output []int
}

var evenTestCases = []testCaseType{
	{
		input:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		output: []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20},
	},
	{
		input:  []int{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
		output: []int{22, 24, 26, 28, 30, 32, 34, 36, 38, 40},
	},
}

func TestEvenNumbers(t *testing.T) {
	for _, testCase := range evenTestCases {
		output := evenNumbers(testCase.input)
		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Wanted %d, got %d", testCase.output, output)
		}
	}
}

var multipleOfFiveTestCases = []testCaseType{
	{
		input:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		output: []int{5, 10, 15, 20},
	},
	{
		input:  []int{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
		output: []int{25, 30, 35, 40},
	},
}

func TestMultipleOfFive(t *testing.T) {
	for _, testCase := range multipleOfFiveTestCases {
		output := multipleOfFive(testCase.input)
		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Wanted %d, got %d", testCase.output, output)
		}
	}
}
