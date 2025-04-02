// write a program to return only the even numbers from this list.
// Sample Input: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
// Sample Output: 2, 4, 6, 8, 10

package main

import (
	"reflect"
	"testing"
)

func TestEvenNumbers(t *testing.T) {
	var sampleInput = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	got := evenNumbers(sampleInput)
	want := []int{2, 4, 6, 8, 10}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

type evenTest struct {
	inputNumbers  []int
	outputNumbers []int
}

var evenTests = []evenTest{
	{inputNumbers: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, outputNumbers: []int{2, 4, 6, 8, 10}},
	{inputNumbers: []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, outputNumbers: []int{12, 14, 16, 18, 20}},
}

func TestEvenNumbers2(t *testing.T) {
	for _, test := range evenTests {
		if !reflect.DeepEqual(evenNumbers(test.inputNumbers), test.outputNumbers) {
			t.Errorf("got %d, wanted %d", test.inputNumbers, test.outputNumbers)
		}
	}

}
