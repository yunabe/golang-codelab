package unittest

// This is a simple example of unit tests in Go.
// https://golang.org/pkg/testing/
// The file of unit tests must have _test.go suffix.
// To run this unit test, execute:
// * go test github.com/yunabe/golang-codelab/unittest
// * go test github.com/yunabe/golang-codelab/unittest -run ^TestIntSumStructured$ -v
// * To learn other options, read the link above.

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestIntSum(t *testing.T) {
	// golang does not support `assert` functions in its unit test library.
	var sum int
	sum = IntSum(1, 2, 3)
	if sum != 6 {
		t.Errorf("Expected 6 but got %d", sum)
	}
}

func TestIntSumStructured(t *testing.T) {
	// Unlike other languages, golang encourages you to write structured unit tests.
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

var (
	goldenFileRe = regexp.MustCompile(`golden\d+\.txt`)
)

func TestGoldenFiles(t *testing.T) {
	suffix := "github.com/yunabe/golang-codelab/unittest"
	cur, err := os.Getwd()
	if err != nil {
		t.Error(err)
		return
	}
	// Confirm the working directory is the directory of this file.
	if !strings.HasSuffix(cur, suffix) {
		t.Errorf("The current dir %q was expected to end with %q but does not.", cur, suffix)
		return
	}
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Error(err)
		return
	}
	for _, file := range files {
		name := file.Name()
		if !goldenFileRe.MatchString(name) {
			continue
		}
		f, err := os.Open(path.Join("testdata", name))
		if err != nil {
			t.Errorf("Failed to open %s: %v", name, err)
			continue
		}
		s := bufio.NewScanner(bufio.NewReader(f))
		lineno := 0
		for s.Scan() {
			lineno++
			line := s.Text()
			pair := strings.Split(line, ":")
			if len(pair) != 2 {
				t.Errorf("The line %d must be split into two with `:` but failed", lineno)
				continue
			}
			want, err := strconv.Atoi(pair[1])
			if err != nil {
				t.Errorf("Failed to parse %q as an integer", pair[1])
				continue
			}
			args := strings.Split(pair[0], ",")
			var intArgs []int
			for _, arg := range args {
				intArg, err := strconv.Atoi(arg)
				if err != nil {
					t.Errorf("Failed to parse %q as an integer", arg)
				} else {
					intArgs = append(intArgs, intArg)
				}
			}
			if len(intArgs) != len(args) {
				continue
			}
			sum := IntSum(intArgs...)
			if sum != want {
				// Don't forget to include the filename in the error message.
				t.Errorf("Expected IntSum(%v...) = %d, but got %d: %s at line %d", intArgs, want, sum, name, lineno)
			}
		}
		if err = f.Close(); err != nil {
			t.Error(err)
		}
	}
}
