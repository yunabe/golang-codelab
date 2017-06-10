package unittest

// This is a simple example of unit tests in Go.
// https://golang.org/pkg/testing/
// The file of unit tests must have _test.go suffix.
// To run this unit test, execute:
// * go test github.com/yunabe/golang-codelab/unittest
// * go test github.com/yunabe/golang-codelab/unittest -run ^TestIntSumStructured$ -v
// * Run other options, read the link above.

import "testing"

func IntSum(args ...int) int {
	sum := 0
	for _, arg := range args {
		sum += arg
	}
	return sum
}

func TestIntSum(t *testing.T) {
	// golang does not support `assert` functions in its unit test library.
	var sum int
	sum = IntSum(1, 2, 3)
	if sum != 6 {
		t.Errorf("Expected 6 but got %d", sum)
	}
}

func TestIntSumStructured(t *testing.T) {
	// Unlike other languages, golang encourage you to write structured uni tests.
	tests := []struct {
		input []int
		want  int
	}{
		{input: nil, want: 0},
		{input: []int{1, 2, 3}, want: 6},
		{input: []int{5, 6, 7, 8}, want: 26},
	}
	for _, test := range tests {
		// The error message is very important.
		if sum := IntSum(test.input...); sum != test.want {
			t.Errorf("Expected IntSum(%v...) = %d, but got %d", test.input, test.want, sum)
		}
	}
}
