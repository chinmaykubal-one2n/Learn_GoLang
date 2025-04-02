// return only the odd numbers from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
// Sample Output: 1, 3, 5, 7, 9
package main

import (
	"reflect"
	"testing"
)

type oddNoType struct {
	input  []int
	output []int
}

var testCases = []oddNoType{
	{
		input:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		output: []int{1, 3, 5, 7, 9},
	},
	{
		input:  []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		output: []int{11, 13, 15, 17, 19},
	},
}

func TestOddNumbers(t *testing.T) {
	for _, test := range testCases {
		wanted := oddNumbers(test.input)
		if !reflect.DeepEqual(wanted, test.output) {
			t.Errorf("wanted %d, got %d", wanted, test.output)
		}
	}
}
