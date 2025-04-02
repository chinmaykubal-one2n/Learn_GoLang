// return only the prime numbers from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
// Sample Output: 2, 3, 5, 7
package main

import (
	"reflect"
	"testing"
)

type primeType struct {
	input  int
	output bool
}

var testCases = []primeType{
	{input: 2, output: true},
	{input: 3, output: true},
	{input: 5, output: true},
	{input: 7, output: true},
}

func TestIsPrime(t *testing.T) {
	for _, test := range testCases {
		result := isPrime(test.input)
		if !reflect.DeepEqual(result, test.output) {
			t.Errorf("Wanted %v, got %v", test.output, result)
		}
	}
}
