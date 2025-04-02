// return only the odd prime numbers from this list
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
// Sample Output: 3, 5, 7

package main

import (
	"reflect"
	"testing"
)

type oddNumbersType struct {
	input  []int
	output []int
}

var testCases = []oddNumbersType{
	{input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, output: []int{1, 3, 5, 7, 9}},
	{input: []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, output: []int{11, 13, 15, 17, 19}},
}

func TestOddNumbers(t *testing.T) {
	for _, test := range testCases {
		output := oddNumbers(test.input)
		if !reflect.DeepEqual(output, test.output) {
			t.Errorf("wanted %d, got %d", test.output, output)
		}
	}
}

type primeNumberType struct {
	input  int
	output bool
}

var testCasesPrime = []primeNumberType{
	{input: 2, output: true},
	{input: 3, output: true},
	{input: 5, output: true},
	{input: 7, output: true},
}

func TestPrimeNumbers(t *testing.T) {
	for _, test := range testCasesPrime {
		output := primeNumbers(test.input)
		if !reflect.DeepEqual(output, test.output) {
			t.Errorf("wanted %v, got %v", test.output, output)
		}
	}
}
