// return only the odd and multiples of 3 greater than 10 from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20
// Sample Output: 15

package main

import (
	"reflect"
	"testing"
)

type testCaseType struct {
	input  []int
	output []int
}

var oddNumbersTestCases = []testCaseType{
	{
		input:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		output: []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19},
	},
	{
		input:  []int{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
		output: []int{21, 23, 25, 27, 29, 31, 33, 35, 37, 39},
	},
}

func TestGetOddNumbers(t *testing.T) {
	for _, testCase := range oddNumbersTestCases {
		output := getOddNumbers(testCase.input)
		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Expected %d, got %d", testCase.output, output)
		}
	}
}

var multiplesOfThreeTestCases = []testCaseType{
	{
		input:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		output: []int{3, 6, 9, 12, 15, 18},
	},
	{
		input:  []int{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
		output: []int{21, 24, 27, 30, 33, 36, 39},
	},
}

func TestGetMultiplesOfThree(t *testing.T) {
	for _, testCase := range multiplesOfThreeTestCases {
		output := getMultiplesOfThree(testCase.input)
		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Expected %d, got %d", testCase.output, output)
		}
	}
}

var greaterThanTenTestCases = []testCaseType{
	{
		input:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		output: []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	},
	{
		input:  []int{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
		output: []int{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
	},
}

func TestGetGreaterThanTen(t *testing.T) {
	for _, testCase := range greaterThanTenTestCases {
		output := getGreaterThanTen(testCase.input)
		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Expected %d, got %d", testCase.output, output)
		}
	}
}
